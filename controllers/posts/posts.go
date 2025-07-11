package posts

import (
	"net/http"

	"go.uber.org/zap"

	"go.opentelemetry.io/otel"
	"ogugu/controllers/common/pgerrors"
	"ogugu/controllers/common/response"
	"ogugu/services/posts"
)

var tracer = otel.Tracer("posts controller")

type PostsController struct {
	log         *zap.Logger
	postService *posts.PostService
}

func New(log *zap.Logger, ps *posts.PostService) *PostsController {
	return &PostsController{
		log:         log,
		postService: ps,
	}
}

func (pc *PostsController) Fetch(w http.ResponseWriter, r *http.Request) {
	spanctx, span := tracer.Start(r.Context(), "fetch all rss")
	defer span.End()

	feed, err := pc.postService.Fetch(spanctx)
	if err != nil {
		pc.log.Error("An error occured while fetching all rss entries", zap.Error(err))
		status, message := pgerrors.Details(err)
		response.Error(w, message, status, pc.log)
		return
	}

	message := "Resoupces Found"
	if len(feed) < 1 {
		message = "No resoupces found"
	}

	response.Success(w, message, http.StatusOK, feed, pc.log)
}
