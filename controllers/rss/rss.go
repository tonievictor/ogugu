package rss

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/oklog/ulid/v2"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"

	"ogugu/controllers/common/pgerrors"
	"ogugu/controllers/common/response"
	"ogugu/models"
	"ogugu/services/rss"
)

var (
	tracer   = otel.Tracer("Rss Controller")
	Validate = validator.New()
)

type RssController struct {
	log        *zap.Logger
	rssService *rss.RssService
}

func New(l *zap.Logger, r *rss.RssService) *RssController {
	return &RssController{
		log:        l,
		rssService: r,
	}
}

// Fetch
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

	feed, err := rc.rssService.Fetch(spanctx)
	if err != nil {
		rc.log.Error("An error occured while fetching all rss entries", zap.Error(err))
		status, message := pgerrors.Details(err)
		response.Error(w, message, status, rc.log)
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
	_, err := rc.rssService.FindByID(spanctx, id)
	if err != nil {
		rc.log.Error("An error occured while deleting an rss entry", zap.String("id", id), zap.Error(err))
		status, message := pgerrors.Details(err)
		response.Error(w, message, status, rc.log)
		return
	}

	err = rc.rssService.DeleteByID(spanctx, id)
	if err != nil {
		rc.log.Error("An error occured while deleting an rss entry", zap.String("id", id), zap.Error(err))
		status, message := pgerrors.Details(err)
		response.Error(w, message, status, rc.log)
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
	feed, err := rc.rssService.FindByID(spanctx, id)
	if err != nil {
		rc.log.Error("resource not found", zap.String("id", id), zap.Error(err))
		status, message := pgerrors.Details(err)
		response.Error(w, message, status, rc.log)
		return
	}

	response.Success(w, "Resource found", http.StatusOK, feed, rc.log)
}

// CreateRss godoc
// @Summary Create a new RSS feed
// @Description Create a new RSS feed by providing the feed's name and link.
// @Tags RSS
// @Accept  json
// @Produce  json
// @Param body body models.CreateRssBody true "Create a new RSS feed"
// @Success 201 {object} response.RssFeed "RSS Feed created"
// @Failure 400 {object} response.Response "Invalid request body"
// @Failure 500 {object} response.Response "Unable to create feed"
// @Router /feed [post]
func (rc *RssController) CreateRss(w http.ResponseWriter, r *http.Request) {
	spanctx, span := tracer.Start(r.Context(), "Create RssFeed")
	defer span.End()

	if r.Body == nil {
		rc.log.Error("request body is missing")
		response.Error(w, "Request body missing", http.StatusBadRequest, rc.log)
		return
	}

	var body models.CreateRssBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		rc.log.Error("invalid request body", zap.Error(err))
		response.Error(w, "Incorrect or Malformed request body", http.StatusBadRequest, rc.log)
		return
	}

	if err = Validate.Struct(body); err != nil {
		rc.log.Error("request body failed some validations", zap.Error(err))
		response.Error(w, err.Error(), http.StatusBadRequest, rc.log)
		return
	}

	id := ulid.Make().String()
	feed, err := rc.rssService.Create(spanctx, id, body)
	if err != nil {
		rc.log.Error("unable to create new rss entry", zap.Error(err))
		status, message := pgerrors.Details(err)
		response.Error(w, message, status, rc.log)
		return
	}

	response.Success(w, "Rss Feed created", http.StatusCreated, feed, rc.log)
}
