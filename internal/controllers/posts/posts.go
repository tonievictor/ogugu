package posts

import (
	"database/sql"
	"errors"
	"net/http"

	"go.uber.org/zap"

	"go.opentelemetry.io/otel"
	"ogugu/internal/controllers/common/response"
	"ogugu/internal/repository/posts"
)

var tracer = otel.Tracer("posts controller")

type Controller struct {
	log      *zap.Logger
	postRepo *posts.Repository
}

func New(log *zap.Logger, ps *posts.Repository) *Controller {
	return &Controller{
		log:      log,
		postRepo: ps,
	}
}

// @Summary		get all posts
// @Description	get all posts
// @Tags			posts
// @Produce		json
// @Success		200		{object}	response.Posts		"Posts found"
// @Failure		default	{object}	response.Response	"Unable to get posts"
// @Router			/posts [get]
func (c *Controller) FetchPosts(w http.ResponseWriter, r *http.Request) {
	spanctx, span := tracer.Start(r.Context(), "fetch all posts")
	defer span.End()

	feed, err := c.postRepo.Fetch(spanctx)
	if err != nil {
		c.log.Error("An error occured while fetching all post entries", zap.Error(err))
		response.Error(w, "internal server error", http.StatusInternalServerError, c.log)
		return
	}

	message := "Resources Found"
	if len(feed) < 1 {
		message = "No resource found"
	}

	response.Success(w, message, http.StatusOK, feed, c.log)
}

// @Summary		get a post
// @Description	get a post by ID
// @Tags			posts
// @Produce		json
// @Param			id		path		string				true	"Post ID"
// @Success		200		{object}	response.Post		"Post with ID found"
// @Failure		404		{object}	response.Response	"Post with ID not found"
// @Failure		default	{object}	response.Response	"Unable to get post with id"
// @Router			/posts/{id} [get]
func (c *Controller) GetPostByID(w http.ResponseWriter, r *http.Request) {
	spanctx, span := tracer.Start(r.Context(), "get a post by id")
	defer span.End()

	id := r.PathValue("id")
	post, err := c.postRepo.GetByID(spanctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.Error(w, "post with id not found", http.StatusNotFound, c.log)
			return
		}
		c.log.Error("unable to get post", zap.Error(err))
		response.Error(w, "internal server error", http.StatusInternalServerError, c.log)
		return
	}

	response.Success(w, "resource with id found", http.StatusOK, post, c.log)
}
