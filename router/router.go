package router

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/redis/go-redis/v9"
	"github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"

	authcontroller "ogugu/controllers/auth"
	postcontroller "ogugu/controllers/posts"
	rsscontroller "ogugu/controllers/rss"
	authservice "ogugu/services/auth"
	postservice "ogugu/services/posts"
	rssservice "ogugu/services/rss"
	userservice "ogugu/services/users"
)

func Routes(db *sql.DB, cache *redis.Client, logger *zap.Logger) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	v1 := chi.NewRouter()
	v1.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte("Up and ready to rumble!!!\n"))
	})

	v1.Get("/swagger/*", httpSwagger.Handler())

	rc := rsscontroller.New(logger, rssservice.New(db))
	v1.Post("/feed", rc.CreateRss)
	v1.Get("/feed/{id}", rc.FindRssByID)
	v1.Get("/feed", rc.Fetch)
	v1.Delete("/feed/{id}", rc.DeleteRssByID)

	ac := authcontroller.New(cache, logger, userservice.New(db), authservice.New(db))
	v1.Post("/signup", ac.Signup)
	v1.Post("/signin", ac.Signin)

	pc := postcontroller.New(logger, postservice.New(db))
	v1.Get("/posts", pc.FetchPosts)
	v1.Get("/posts/{id}", pc.GetPostByID)

	r.Mount("/v1", v1)
	return r
}
