package auth

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/oklog/ulid/v2"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"ogugu/controllers/common/pgerrors"
	"ogugu/controllers/common/response"
	"ogugu/models"
	"ogugu/services/auth"
	"ogugu/services/users"
)

var (
	tracer   = otel.Tracer("Auth Controller")
	Validate = validator.New()
)

type AuthController struct {
	cache       *redis.Client
	log         *zap.Logger
	userService *users.UserService
	authService *auth.AuthService
}

func New(c *redis.Client, l *zap.Logger, u *users.UserService, a *auth.AuthService) *AuthController {
	return &AuthController{
		cache:       c,
		log:         l,
		userService: u,
		authService: a,
	}
}

// @Summary			 sign up
// @Description  create a new account
// @Tags         account
// @Accept       json
// @Produce      json
// @Param body body models.CreateUserBody true "body"
// @Success 201  {object} response.RssFeed
// @Failure 400  {object} response.Response
// @Failure 500  {object} response.Response
// @Router /signup [post]
func (ac *AuthController) Signup(w http.ResponseWriter, r *http.Request) {
	spanctx, span := tracer.Start(r.Context(), "Sign Up")
	defer span.End()

	if r.Body == nil {
		ac.log.Error("request body is missing")
		response.Error(w, "Request body missing", http.StatusBadRequest, ac.log)
		return
	}

	var body models.CreateUserBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		ac.log.Error("invalid request body", zap.Error(err))
		response.Error(w, "Incorrect or Malformed request body", http.StatusBadRequest, ac.log)
		return
	}

	if err = Validate.Struct(body); err != nil {
		ac.log.Error("invalid request body", zap.Error(err))
		response.Error(w, err.Error(), http.StatusBadRequest, ac.log)
		return
	}

	id := ulid.Make().String()
	user, err := ac.userService.CreateUser(spanctx, id, body)
	if err != nil {
		ac.log.Error("cannot create new user", zap.Error(err))
		status, message := pgerrors.Details(err)
		response.Error(w, message, status, ac.log)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(body.Password), 4)
	if err != nil {
		ac.log.Error("cannot create new user", zap.Error(err))
		response.Error(w, "An error occured while creating the user", http.StatusInternalServerError, ac.log)
		return
	}

	err = ac.authService.CreateAuth(spanctx, user.ID, string(hashed))
	if err != nil {
		ac.log.Error("cannot create user auth details", zap.Error(err))
		status, message := pgerrors.Details(err)
		response.Error(w, message, status, ac.log)
		return
	}

	response.Success(w, "Sign up successfull", http.StatusCreated, user, ac.log)
}
