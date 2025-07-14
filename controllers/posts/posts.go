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

// @Summary      get all posts
// @Description  get all posts
// @Tags         posts
// @Produce      json
// @Success 200 {object} response.Posts "Posts found"
// @Failure default {object} response.Response "Unable to get posts"
// @Router /posts [get]
func (pc *PostsController) FetchPosts(w http.ResponseWriter, r *http.Request) {
	spanctx, span := tracer.Start(r.Context(), "fetch all posts")
	defer span.End()

	feed, err := pc.postService.Fetch(spanctx)
	if err != nil {
		pc.log.Error("An error occured while fetching all post entries", zap.Error(err))
		status, message := pgerrors.Details(err)
		response.Error(w, message, status, pc.log)
		return
	}

	message := "Resources Found"
	if len(feed) < 1 {
		message = "No resource found"
	}

	response.Success(w, message, http.StatusOK, feed, pc.log)
}

// @Summary      get a post
// @Description  get a post by ID
// @Tags         posts
// @Produce      json
// @Param        id   path      string  true  "Post ID"
// @Success 200 {object} response.Post "Post with ID found"
// @Failure 404 {object} response.Response "Post with ID not found"
// @Failure default {object} response.Response "Unable to get post with id"
// @Router /posts/{id} [get]
func (pc *PostsController) GetPostByID(w http.ResponseWriter, r *http.Request) {
	spanctx, span := tracer.Start(r.Context(), "get a post by id")
	defer span.End()

	id := r.PathValue("id")
	post, err := pc.postService.GetByID(spanctx, id)
	if err != nil {
		pc.log.Error("unable to get post", zap.Error(err))
		status, message := pgerrors.Details(err)
		response.Error(w, message, status, pc.log)
		return
	}

	response.Success(w, "resource with id found", http.StatusOK, post, pc.log)
}
