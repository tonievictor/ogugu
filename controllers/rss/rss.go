package rss

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"time"

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
	tracer   = otel.Tracer("rss controller")
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

// @Summary Find all RSS feeds
// @Description Retrieve all RSS Feeds in the database.
// @Tags rss
// @Produce  json
// @Success 200 {object} response.RssFeeds "RSS Feeds found"
// @Failure 400 {object} response.Response "Invalid request"
// @Failure 404 {object} response.Response "RSS Feed not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /feed [get]
func (rc *RssController) Fetch(w http.ResponseWriter, r *http.Request) {
	spanctx, span := tracer.Start(r.Context(), "fetch all rss")
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

// @Summary Delete an RSS feed by its ID
// @Description Delete an existing RSS feed using its unique ID.
// @Tags rss
// @Produce  json
// @Param id path string true "ID of the RSS feed to retrieve"
// @Success 204 {object} response.Response "RSS Feed deleted"
// @Failure 400 {object} response.Response "Invalid request"
// @Failure 404 {object} response.Response "RSS Feed not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /feed/{id} [delete]
func (rc *RssController) DeleteRssByID(w http.ResponseWriter, r *http.Request) {
	spanctx, span := tracer.Start(r.Context(), "delete rss by id")
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

	_, err = rc.rssService.DeleteByID(spanctx, id)
	if err != nil {
		rc.log.Error("An error occured while deleting an rss entry", zap.String("id", id), zap.Error(err))
		status, message := pgerrors.Details(err)
		response.Error(w, message, status, rc.log)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

// @Summary Find an RSS feed by its ID
// @Description Retrieve an existing RSS feed using its unique ID.
// @Tags rss
// @Produce  json
// @Param id path string true "ID of the RSS feed to retrieve"
// @Success 200 {object} response.RssFeed "RSS Feed found"
// @Failure 400 {object} response.Response "Invalid request"
// @Failure 404 {object} response.Response "RSS Feed not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /feed/{id} [get]
func (rc *RssController) FindRssByID(w http.ResponseWriter, r *http.Request) {
	spanctx, span := tracer.Start(r.Context(), "fetch rss by id")
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

// @Summary Create a new RSS feed
// @Description Create a new RSS feed by providing the feed's name and link.
// @Tags rss
// @Accept  json
// @Produce  json
// @Param body body models.CreateRssBody true "Create a new RSS feed"
// @Success 201 {object} response.RssFeed "RSS Feed created"
// @Failure 400 {object} response.Response "Invalid request body"
// @Failure 500 {object} response.Response "Unable to create feed"
// @Router /feed [post]
func (rc *RssController) CreateRss(w http.ResponseWriter, r *http.Request) {
	spanctx, span := tracer.Start(r.Context(), "create rss feed")
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

	meta, msg, err := getRssInfo(body.Link, spanctx)
	if err != nil {
		rc.log.Error(msg, zap.Error(err))
		response.Error(w, msg, http.StatusBadRequest, rc.log)
		return
	}

	id := ulid.Make().String()
	feed, err := rc.rssService.Create(spanctx, id, body.Link, meta)
	if err != nil {
		rc.log.Error("unable to create new rss entry", zap.Error(err))
		status, message := pgerrors.Details(err)
		response.Error(w, message, status, rc.log)
		return
	}

	response.Success(w, "Rss Feed created", http.StatusCreated, feed, rc.log)
}

func getRssInfo(url string, ctx context.Context) (models.RSSMeta, string, error) {
	_, span := tracer.Start(ctx, "get rss metadata")
	defer span.End()
	client := http.Client{
		Timeout: time.Second * 10,
	}

	res, err := client.Get(url)
	if err != nil || res.StatusCode != 200 {
		return models.RSSMeta{}, "could not fetch from rss link", err
	}

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return models.RSSMeta{}, "an error occured while reading rss body", err
	}

	var body models.RSSMeta
	err = xml.Unmarshal(data, &body)
	if err != nil {
		return models.RSSMeta{}, "The URL does not point to a valid RSS feed", err
	}

	return body, "", nil
}
