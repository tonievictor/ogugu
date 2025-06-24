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

// Fetch godoc
// @Summary Find all RSS feeds
// @Description Retrieve all RSS Feeds in the database.
// @Tags RSS
// @Accept  json
// @Produce  json
// @Success 200 {object} response.RssFeeds "RSS Feeds found"
// @Failure 400 {object} response.Response "Invalid request"
// @Failure 404 {object} response.Response "RSS Feed not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /feed [get]
func (rc *RssController) Fetch(w http.ResponseWriter, r *http.Request) {
	spanctx, span := tracer.Start(r.Context(), "Find RssFeedByID")
	defer span.End()

	feed, err := rc.rss.Fetch(spanctx)
	if err != nil {
		response.Error(w, "Resource not found", http.StatusNotFound, err.Error(), rc.log)
		return
	}

	message := "Resources Found"
	if len(feed) < 1 {
		message = "No resources found"
	}

	response.Success(w, message, http.StatusOK, feed, rc.log)
}

// DeleteRssByID godoc
// @Summary Delete an RSS feed by its ID
// @Description Delete an existing RSS feed using its unique ID.
// @Tags RSS
// @Accept  json
// @Produce  json
// @Param id path string true "ID of the RSS feed to retrieve"
// @Success 204 {object} response.Response "RSS Feed deleted"
// @Failure 400 {object} response.Response "Invalid request"
// @Failure 404 {object} response.Response "RSS Feed not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /feed/{id} [delete]
func (rc *RssController) DeleteRssByID(w http.ResponseWriter, r *http.Request) {
	spanctx, span := tracer.Start(r.Context(), "Delete Rss by ID")
	defer span.End()

	id := r.PathValue("id")
	// DeleteByID
	_, err := rc.rss.FindByID(spanctx, id)
	if err != nil {
		response.Error(w, "Resource not found", http.StatusNotFound, err.Error(), rc.log)
		return
	}

	err = rc.rss.DeleteByID(spanctx, id)
	if err != nil {
		msg := "An error occured while deleting the resource"
		response.Error(w, msg, http.StatusInternalServerError, err.Error(), rc.log)
		return
	}

	response.Success(w, "Resource deleted successfully", http.StatusNoContent, "", rc.log)
}

// FindRssByID godoc
// @Summary Find an RSS feed by its ID
// @Description Retrieve an existing RSS feed using its unique ID.
// @Tags RSS
// @Accept  json
// @Produce  json
// @Param id path string true "ID of the RSS feed to retrieve"
// @Success 200 {object} response.RssFeed "RSS Feed found"
// @Failure 400 {object} response.Response "Invalid request"
// @Failure 404 {object} response.Response "RSS Feed not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /feed/{id} [get]
func (rc *RssController) FindRssByID(w http.ResponseWriter, r *http.Request) {
	spanctx, span := tracer.Start(r.Context(), "Find RssFeedByID")
	defer span.End()

	id := r.PathValue("id")
	feed, err := rc.rss.FindByID(spanctx, id)
	if err != nil {
		response.Error(w, "Resource not found", http.StatusNotFound, err.Error(), rc.log)
		return
	}

	response.Success(w, "Resource found", http.StatusOK, feed, rc.log)
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
}
