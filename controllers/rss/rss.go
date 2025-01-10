package rss

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/oklog/ulid/v2"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"

	"ogugu/response"
	feed "ogugu/services/rss"
)

var (
	tracer   = otel.Tracer("Rss Controller")
	Validate = validator.New()
)

type RssController struct {
	cache *redis.Client
	log   *zap.Logger
	rss   *feed.RssService
}

func New(
	cache *redis.Client, log *zap.Logger, rss *feed.RssService,
) *RssController {
	return &RssController{
		cache: cache,
		log:   log,
		rss:   rss,
	}
}

type CreateRssBody struct {
	Name string `json:"name" validate:"required"`
	Link string `json:"link" validate:"required"`
}

// CreateRss godoc
// @Summary Create a new RSS feed
// @Description Create a new RSS feed by providing the feed's name and link.
// @Tags RSS
// @Accept  json
// @Produce  json
// @Param body body CreateRssBody true "Create a new RSS feed"
// @Success 201 {object} response.RssFeed "RSS Feed created"
// @Failure 400 {object} response.Response "Invalid request body"
// @Failure 500 {object} response.Response "Unable to create feed"
// @Router /feed [post]
func (rc *RssController) CreateRss(w http.ResponseWriter, r *http.Request) {
	spanctx, span := tracer.Start(r.Context(), "Create RssFeed")
	defer span.End()

	if r.Body == nil {
		err := errors.New("Request body cannot be empty")
		response.Error(w, "Empty request body", http.StatusBadRequest, err.Error(), rc.log)
		return
	}

	var body CreateRssBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		err := errors.New("The request body is malformed or not valid JSON")
		response.Error(w, "Invalid request body", http.StatusBadRequest, err.Error(), rc.log)
		return
	}

	if err = Validate.Struct(body); err != nil {
		err := err.(validator.ValidationErrors)
		response.Error(w, "Invalid request body", http.StatusBadRequest, err.Error(), rc.log)
		return
	}
	id := ulid.Make().String()
	feed, err := rc.rss.Create(spanctx, body.Name, body.Link, id)
	if err != nil {
		s := http.StatusInternalServerError
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			s = http.StatusConflict
		}
		response.Error(w, "Unable to create feed", s, err.Error(), rc.log)
		return
	}

	response.Success(w, "Rss Feed created", http.StatusCreated, feed, rc.log)
	return
}
