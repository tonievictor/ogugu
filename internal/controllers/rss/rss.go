package rss

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/oklog/ulid/v2"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"

	"ogugu/internal/controllers/common/response"
	"ogugu/internal/models"
	"ogugu/internal/repository/rss"
)

var (
	tracer   = otel.Tracer("rss controller")
	Validate = validator.New()
)

type Controller struct {
	log     *zap.Logger
	rssRepo *rss.Repository
}

func New(l *zap.Logger, r *rss.Repository) *Controller {
	return &Controller{
		log:     l,
		rssRepo: r,
	}
}

// @Summary		Find all RSS feeds
// @Description	Retrieve all RSS Feeds in the database.
// @Tags			rss
// @Produce		json
// @Success		200		{object}	response.RssFeeds	"RSS Feeds found"
// @Failure		400		{object}	response.Response	"Invalid request"
// @Failure		404		{object}	response.Response	"RSS Feed not found"
// @Failure		500		{object}	response.Response	"An error occured on the server"
// @Failure		default	{object}	response.Response	"An error occured"
// @Router			/feed [get]
func (c *Controller) Fetch(w http.ResponseWriter, r *http.Request) {
	spanctx, span := tracer.Start(r.Context(), "fetch all rss")
	defer span.End()

	feed, err := c.rssRepo.Fetch(spanctx)
	if err != nil {
		c.log.Error("An error occured while fetching all rss entries", zap.Error(err))
		response.Error(w, "internal server error", http.StatusInternalServerError, c.log)
		return
	}

	message := "Resources Found"
	if len(feed) < 1 {
		message = "No resources found"
	}

	response.Success(w, message, http.StatusOK, feed, c.log)
}

// @Summary		Delete an RSS feed by its ID
// @Description	Delete an existing RSS feed using its unique ID.
// @Tags			rss
// @Produce		json
// @Param			id		path		string				true	"ID of the RSS feed to retrieve"
// @Success		204		{object}	response.Response	"RSS Feed deleted"
// @Failure		400		{object}	response.Response	"Invalid request"
// @Failure		404		{object}	response.Response	"RSS Feed not found"
// @Failure		500		{object}	response.Response	"An error occured on the server"
// @Failure		default	{object}	response.Response	"An error occured"
// @Router			/feed/{id} [delete]
func (c *Controller) DeleteRssByID(w http.ResponseWriter, r *http.Request) {
	spanctx, span := tracer.Start(r.Context(), "delete rss by id")
	defer span.End()

	id := r.PathValue("id")
	_, err := c.rssRepo.FindByID(spanctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.log.Warn("rss entry not found", zap.String("id", id))
			response.Error(w, "rss with id not found", http.StatusNotFound, c.log)
			return
		}

		c.log.Error("An error occured while deleting an rss entry", zap.String("id", id), zap.Error(err))
		response.Error(w, "internal server error", http.StatusInternalServerError, c.log)
		return
	}

	_, err = c.rssRepo.DeleteByID(spanctx, id)
	if err != nil {
		c.log.Error("An error occured while deleting an rss entry", zap.String("id", id), zap.Error(err))
		response.Error(w, "internal server error", http.StatusInternalServerError, c.log)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

// @Summary		Find an RSS feed by its ID
// @Description	Retrieve an existing RSS feed using its unique ID.
// @Tags			rss
// @Produce		json
// @Param			id		path		string				true	"ID of the RSS feed to retrieve"
// @Success		200		{object}	response.RssFeed	"RSS Feed found"
// @Failure		400		{object}	response.Response	"Invalid or malformed request body"
// @Failure		404		{object}	response.Response	"RSS Feed not found"
// @Failure		500		{object}	response.Response	"An error occured on the server"
// @Failure		default	{object}	response.Response	"An error occured"
// @Router			/feed/{id} [get]
func (c *Controller) FindRssByID(w http.ResponseWriter, r *http.Request) {
	spanctx, span := tracer.Start(r.Context(), "fetch rss by id")
	defer span.End()

	id := r.PathValue("id")
	feed, err := c.rssRepo.FindByID(spanctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.log.Warn("rss entry not found", zap.String("id", id))
			response.Error(w, "rss with id not found", http.StatusNotFound, c.log)
			return
		}

		c.log.Error("an error occured while fetching rss entry", zap.String("id", id), zap.Error(err))
		response.Error(w, "internal server error", http.StatusInternalServerError, c.log)
		return
	}

	response.Success(w, "Resource found", http.StatusOK, feed, c.log)
}

// @Summary		Create a new RSS feed
// @Description	Create a new RSS feed by providing the feed's name and link.
// @Tags			rss
// @Accept			json
// @Produce		json
// @Param			body	body		models.CreateRssBody	true	"Create a new RSS feed"
// @Success		201		{object}	response.RssFeed		"RSS Feed created"
// @Failure		400		{object}	response.Response		"Invalid or malformed request body"
// @Failure		500		{object}	response.Response		"An error occured on the server"
// @Failure		default	{object}	response.Response		"An error occured"
// @Router			/feed [post]
func (c *Controller) CreateRss(w http.ResponseWriter, r *http.Request) {
	spanctx, span := tracer.Start(r.Context(), "create rss feed")
	defer span.End()

	if r.Body == nil {
		c.log.Error("request body is missing")
		response.Error(w, "Request body missing", http.StatusBadRequest, c.log)
		return
	}

	var body models.CreateRssBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		c.log.Error("invalid request body", zap.Error(err))
		response.Error(w, "Incorrect or Malformed request body", http.StatusBadRequest, c.log)
		return
	}

	if err = Validate.Struct(body); err != nil {
		c.log.Error("request body failed some validations", zap.Error(err))
		response.Error(w, err.Error(), http.StatusBadRequest, c.log)
		return
	}

	meta, err := getRSSMeta(body.Link)
	if err != nil {
		c.log.Error(err.Error(), zap.Error(err))
		response.Error(w, "an error occured while fetching rss metadata", http.StatusUnprocessableEntity, c.log)
		return
	}

	if meta.Channel.Title == "" {
		meta.Channel.Title = "Untitled Feed"
	}

	id := ulid.Make().String()
	feed, err := c.rssRepo.Create(spanctx, id, body.Link, meta)
	if err != nil {
		c.log.Error("could not create new feed", zap.Error(err))
		response.Error(w, "could not create new feed", http.StatusInternalServerError, c.log)
		return
	}

	response.Success(w, "rss feed created successfully", http.StatusCreated, feed, c.log)
}

func getRSSMeta(link string) (models.RSSMeta, error) {
	res, err := http.Get(link)
	if err != nil {
		return models.RSSMeta{}, err
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return models.RSSMeta{}, err
	}

	var meta models.RSSMeta
	err = xml.Unmarshal(body, &meta)
	if err != nil {
		return models.RSSMeta{}, err
	}

	lastModified := res.Header.Get("Last-Modified")
	if lastModified == "" {
		meta.Channel.LastModified = time.Now().Format(time.RFC1123)
	} else {
		meta.Channel.LastModified = lastModified
	}
	return meta, nil
}
