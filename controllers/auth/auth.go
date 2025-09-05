package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/oklog/ulid/v2"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"ogugu/controllers/common/pgerrors"
	"ogugu/controllers/common/response"
	"ogugu/models"
	"ogugu/repository/auth"
	"ogugu/repository/users"
)

var (
	tracer   = otel.Tracer("Auth Controller")
	Validate = validator.New()
)

type Controller struct {
	cache    *redis.Client
	log      *zap.Logger
	userRepo *users.Repository
	authRepo *auth.Repository
}

func New(c *redis.Client, l *zap.Logger, u *users.Repository, a *auth.Repository) *Controller {
	return &Controller{
		cache:    c,
		log:      l,
		userRepo: u,
		authRepo: a,
	}
}

// @Summary		sign out
// @Description	sign out from current session
// @Security		BearerAuth
// @Tags			account
// @Accept			json
// @Produce		json
// @Sucess			204
// @Failure		401		{object}	response.Response
// @Failure		500		{object}	response.Response
// @Failure		default	{object}	response.Response
// @Router			/signout [delete]
func (c *Controller) Signout(w http.ResponseWriter, r *http.Request) {
	spanctx, span := tracer.Start(r.Context(), "sign out")
	defer span.End()

	sess := r.Context().Value(models.AuthSessionKey).(models.Session)

	err := c.cache.SetEx(spanctx, sess.ID, "", time.Second).Err()
	if err != nil {
		c.log.Error("could not delete user session", zap.Error(err))
		response.Error(w, "An error occured while deleting the session", http.StatusInternalServerError, c.log)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

// @Summary		sign in
// @Description	signin to an existing account
// @Tags			account
// @Accept			json
// @Produce		json
// @Param			body	body		models.SigninBody	true	"body"
// @Success		200		{object}	response.UserWithAuth
// @Failure		400		{object}	response.Response
// @Failure		500		{object}	response.Response
// @Router			/signin [post]
func (c *Controller) Signin(w http.ResponseWriter, r *http.Request) {
	spanctx, span := tracer.Start(r.Context(), "Sign in")
	defer span.End()

	if r.Body == nil {
		c.log.Error("request body is missing")
		response.Error(w, "Request body missing", http.StatusBadRequest, c.log)
		return
	}

	var body models.SigninBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		c.log.Error("invalid request body", zap.Error(err))
		response.Error(w, "Incorrect or Malformed request body", http.StatusBadRequest, c.log)
		return
	}

	if err = Validate.Struct(body); err != nil {
		c.log.Error("invalid request body", zap.Error(err))
		response.Error(w, err.Error(), http.StatusBadRequest, c.log)
		return
	}

	user, err := c.userRepo.GetUser(spanctx, "email", body.Email)
	if err != nil {
		c.log.Error("an error occured while fetching user", zap.Error(err))
		status, _ := pgerrors.Details(err)
		response.Error(w, "Login Failed, check credentials", status, c.log)
		return
	}

	hashpwd, err := c.authRepo.GetPasswordWithUserID(spanctx, user.ID)
	if err != nil {
		c.log.Error("an error occured while fetching user", zap.Error(err))
		status, _ := pgerrors.Details(err)
		response.Error(w, "Login Failed, check credentials", status, c.log)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashpwd), []byte(body.Password))
	if err != nil {
		c.log.Error("password mismatch", zap.Error(err))
		status, _ := pgerrors.Details(err)
		response.Error(w, "Login Failed, check credentials", status, c.log)
		return
	}

	sessionid := ulid.Make().String()
	session, err := json.Marshal(models.Session{
		ID:         sessionid,
		UserID:     user.ID,
		CreatedAt:  time.Now(),
		ExpiryTime: time.Now().Add(time.Hour * 24 * 3),
	})
	if err != nil {
		c.log.Error("unable to marshal session", zap.Error(err))
		response.Error(w, "Login Failed, please try again", http.StatusInternalServerError, c.log)
		return
	}

	err = c.cache.Set(spanctx, sessionid, session, time.Second*259200).Err()
	if err != nil && err != redis.Nil {
		c.log.Error("unable to create session", zap.Error(err))
		response.Error(w, "Login Failed, please try again", http.StatusInternalServerError, c.log)
		return
	}

	data := models.UserWithAuth{
		User:      user,
		AuthToken: sessionid,
	}
	response.Success(w, "Login Successful", http.StatusOK, data, c.log)
}

// @Summary		sign up
// @Description	create a new account
// @Tags			account
// @Accept			json
// @Produce		json
// @Param			body	body		models.CreateUserBody	true	"body"
// @Success		201		{object}	response.User
// @Failure		400		{object}	response.Response
// @Failure		500		{object}	response.Response
// @Router			/signup [post]
func (c *Controller) Signup(w http.ResponseWriter, r *http.Request) {
	spanctx, span := tracer.Start(r.Context(), "Sign Up")
	defer span.End()

	if r.Body == nil {
		c.log.Error("request body is missing")
		response.Error(w, "Request body missing", http.StatusBadRequest, c.log)
		return
	}

	var body models.CreateUserBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		c.log.Error("invalid request body", zap.Error(err))
		response.Error(w, "Incorrect or Malformed request body", http.StatusBadRequest, c.log)
		return
	}

	if err = Validate.Struct(body); err != nil {
		c.log.Error("invalid request body", zap.Error(err))
		response.Error(w, err.Error(), http.StatusBadRequest, c.log)
		return
	}

	id := ulid.Make().String()
	user, err := c.userRepo.CreateUser(spanctx, id, body)
	if err != nil {
		c.log.Error("cannot create new user", zap.Error(err))
		status, message := pgerrors.Details(err)
		response.Error(w, message, status, c.log)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(body.Password), 4)
	if err != nil {
		c.log.Error("cannot create new user", zap.Error(err))
		response.Error(w, "An error occured while creating the user", http.StatusInternalServerError, c.log)
		return
	}

	err = c.authRepo.CreateAuth(spanctx, user.ID, string(hashed))
	if err != nil {
		c.log.Error("cannot create user auth details", zap.Error(err))
		status, message := pgerrors.Details(err)
		response.Error(w, message, status, c.log)
		return
	}

	response.Success(w, "Sign up successfull", http.StatusCreated, user, c.log)
}
