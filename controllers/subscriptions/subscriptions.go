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

type SubscriptionController struct {
	cache      *redis.Client
	log        *zap.Logger
	subservice *subscriptions.SubscriptionService
}

func New(cache *redis.Client,
	log *zap.Logger,
	ss *subscriptions.SubscriptionService,
) *SubscriptionController {
	return &SubscriptionController{
		cache:      cache,
		log:        log,
		subservice: ss,
	}
}

// @Summary			 subscribe
// @Description  subscribe to an rss feed
// @Tags         subscription
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param body body models.SubscriptionBody true "body"
// @Success 201  {object} response.Subscription
// @Failure 400  {object} response.Response
// @Failure 500  {object} response.Response
// @Failure default  {object} response.Response
// @Router /subscriptions [post]
func (sc *SubscriptionController) Subscribe(w http.ResponseWriter, r *http.Request) {
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

	u := r.Context().Value(models.AuthSession).(models.Session)
	id := ulid.Make().String()
	sub, err := sc.subservice.CreateSub(spanctx, id, u.UserID, body.RssID)
	if err != nil {
		sc.log.Error("could not add subscription", zap.Error(err))
		status, message := pgerrors.Details(err)
		response.Error(w, message, status, sc.log)
		return
	}

	response.Success(w, "subscription successful", http.StatusCreated, sub, sc.log)
}

// @Summary			 unsubscribe
// @Description  unsubscribe from an rss feed
// @Tags         subscription
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param body body models.SubscriptionBody true "body"
// @Success 204
// @Failure 400  {object} response.Response
// @Failure 500  {object} response.Response
// @Failure default  {object} response.Response
// @Router /subscriptions [delete]
func (sc *SubscriptionController) Unsubscribe(w http.ResponseWriter, r *http.Request) {
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

	u := r.Context().Value(models.AuthSession).(models.Session)
	_, err = sc.subservice.DeleteSub(spanctx, u.UserID, body.RssID)
	if err != nil {
		sc.log.Error("could not delete subscription", zap.Error(err))
		status, message := pgerrors.Details(err)
		response.Error(w, message, status, sc.log)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
