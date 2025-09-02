package subscriptions

import (
	"encoding/json"
	"net/http"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"ogugu/controllers/common/pgerrors"
	"ogugu/controllers/common/response"
	"ogugu/models"

	"github.com/go-playground/validator/v10"
	"github.com/oklog/ulid/v2"
	"ogugu/services/subscriptions"
)

var (
	tracer   = otel.Tracer("subscriptions controller")
	Validate = validator.New()
)

type Controller struct {
	cache   *redis.Client
	log     *zap.Logger
	service *subscriptions.Service
}

func New(cache *redis.Client,
	log *zap.Logger,
	ss *subscriptions.Service,
) *Controller {
	return &Controller{
		cache:   cache,
		log:     log,
		service: ss,
	}
}

// @Summary		get subscriptions
// @Description	get current user's subscriptions
// @Tags			subscription
// @Security		BearerAuth
// @Accept			json
// @Produce		json
// @Success		200		{object}	response.Subscription
// @Failure		400		{object}	response.Response
// @Failure		500		{object}	response.Response
// @Failure		default	{object}	response.Response
// @Router			/subscriptions [get]
func (sc *Controller) GetUserSubs(w http.ResponseWriter, r *http.Request) {
	spanctx, span := tracer.Start(r.Context(), "get user's subscriptions")
	defer span.End()

	session := r.Context().Value(models.AuthSessionKey).(models.Session)

	subs, err := sc.service.GetSubsByUserID(spanctx, session.UserID)
	if err != nil {
		sc.log.Error("could not get user's subscriptions", zap.Error(err))
		status, message := pgerrors.Details(err)
		response.Error(w, message, status, sc.log)
		return
	}

	msg := "Resources found"
	if len(subs) == 0 {
		msg = "No resource"
	}
	response.Success(w, msg, http.StatusOK, subs, sc.log)
}

// @Summary		subscribe
// @Description	subscribe to an rss feed
// @Tags			subscription
// @Security		BearerAuth
// @Accept			json
// @Produce		json
// @Param			body	body		models.SubscriptionBody	true	"body"
// @Success		201		{object}	response.Subscription
// @Failure		400		{object}	response.Response
// @Failure		500		{object}	response.Response
// @Failure		default	{object}	response.Response
// @Router			/subscriptions [post]
func (sc *Controller) Subscribe(w http.ResponseWriter, r *http.Request) {
	spanctx, span := tracer.Start(r.Context(), "subscribe to an rss feed")
	defer span.End()

	if r.Body == nil {
		sc.log.Error("request body is missing")
		response.Error(w, "Request body missing", http.StatusBadRequest, sc.log)
		return
	}

	var body models.SubscriptionBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		sc.log.Error("Could not read request body", zap.Error(err))
		response.Error(w, "Unable to read request body", http.StatusBadRequest, sc.log)
		return
	}

	if err = Validate.Struct(body); err != nil {
		sc.log.Error("request body failed some validations", zap.Error(err))
		response.Error(w, err.Error(), http.StatusBadRequest, sc.log)
		return
	}

	u := r.Context().Value(models.AuthSessionKey).(models.Session)
	id := ulid.Make().String()
	sub, err := sc.service.CreateSub(spanctx, id, u.UserID, body.RssID)
	if err != nil {
		sc.log.Error("could not add subscription", zap.Error(err))
		status, message := pgerrors.Details(err)
		response.Error(w, message, status, sc.log)
		return
	}

	response.Success(w, "subscription successful", http.StatusCreated, sub, sc.log)
}

// @Summary		unsubscribe
// @Description	unsubscribe from an rss feed
// @Tags			subscription
// @Security		BearerAuth
// @Accept			json
// @Produce		json
// @Param			body	body	models.SubscriptionBody	true	"body"
// @Success		204
// @Failure		400		{object}	response.Response
// @Failure		500		{object}	response.Response
// @Failure		default	{object}	response.Response
// @Router			/subscriptions [delete]
func (sc *Controller) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	spanctx, span := tracer.Start(r.Context(), "unsubscribe from an rss feed")
	defer span.End()

	if r.Body == nil {
		sc.log.Error("request body is missing")
		response.Error(w, "Request body missing", http.StatusBadRequest, sc.log)
		return
	}

	var body models.SubscriptionBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		sc.log.Error("Could not read request body", zap.Error(err))
		response.Error(w, "Unable to read request body", http.StatusBadRequest, sc.log)
		return
	}

	if err = Validate.Struct(body); err != nil {
		sc.log.Error("request body failed some validations", zap.Error(err))
		response.Error(w, err.Error(), http.StatusBadRequest, sc.log)
		return
	}

	session := r.Context().Value(models.AuthSessionKey).(models.Session)
	_, err = sc.service.DeleteSub(spanctx, session.UserID, body.RssID)
	if err != nil {
		sc.log.Error("could not delete subscription", zap.Error(err))
		status, message := pgerrors.Details(err)
		response.Error(w, message, status, sc.log)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

// @Summary		get posts
// @Description	get posts from feed that user is subscribed to
// @Tags			subscription
// @Security		BearerAuth
// @Accept			json
// @Produce		json
// @Success		200		{object}	response.FeedPosts
// @Failure		400		{object}	response.Response
// @Failure		500		{object}	response.Response
// @Failure		default	{object}	response.Response
// @Router			/subscriptions/posts [get]
func (sc *Controller) GetPostFromSub(w http.ResponseWriter, r *http.Request) {
	spanctx, span := tracer.Start(r.Context(), "get post from sub")
	defer span.End()

	session := r.Context().Value(models.AuthSessionKey).(models.Session)
	posts, err := sc.service.GetPostFromSubScriptions(spanctx, session.UserID)
	if err != nil {
		sc.log.Error("cannot get posts", zap.Error(err))
		status, msg := pgerrors.Details(err)
		response.Error(w, msg, status, sc.log)
		return
	}

	msg := "resources found"
	if len(posts) == 0 {
		msg = "no resource found"
	}
	response.Success(w, msg, http.StatusOK, posts, sc.log)
}
