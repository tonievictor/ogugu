package controllers

import (
	"net/http"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"

	"ogugu/services"
)

var tracer = otel.Tracer("Rss Controller")

type RssController struct {
	cache *redis.Client
	log   *zap.Logger
	rss   *services.RssService
}

func NewRssController(cache *redis.Client, log *zap.Logger, rss *services.RssService) *RssController {
	return &RssController{
		cache: cache,
		log:   log,
		rss:   rss,
	}
}

func (rc *RssController) CreateRss(w http.ResponseWriter, r *http.Request) {}
