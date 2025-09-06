package router

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/redis/go-redis/v9"
	"github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"

	authcontroller "ogugu/internal/controllers/auth"
	postcontroller "ogugu/internal/controllers/posts"
	rsscontroller "ogugu/internal/controllers/rss"
	subcontroller "ogugu/internal/controllers/subscriptions"
	authRepo "ogugu/internal/repository/auth"
	postRepo "ogugu/internal/repository/posts"
	rssRepo "ogugu/internal/repository/rss"
	subRepo "ogugu/internal/repository/subscriptions"
	userRepo "ogugu/internal/repository/users"
)

func Routes(db *sql.DB, cache *redis.Client, logger *zap.Logger) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1 := chi.NewRouter()
	v1.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte("Up and ready to rumble!!!\n"))
	})

	v1.Get("/swagger/*", httpSwagger.Handler())

	rc := rsscontroller.New(logger, rssRepo.New(db))
	v1.Post("/feed", rc.CreateRss)
	v1.Get("/feed/{id}", rc.FindRssByID)
	v1.Get("/feed", rc.Fetch)
	v1.Delete("/feed/{id}", rc.DeleteRssByID)

	ac := authcontroller.New(cache, logger, userRepo.New(db), authRepo.New(db))
	v1.Post("/signup", ac.Signup)
	v1.Post("/signin", ac.Signin)
	v1.Delete("/signout", IsAuthenticated(cache, logger, ac.Signout))

	pc := postcontroller.New(logger, postRepo.New(db))
	v1.Get("/posts", pc.FetchPosts)
	v1.Get("/posts/{id}", pc.GetPostByID)

	sc := subcontroller.New(cache, logger, subRepo.New(db))
	v1.Post("/subscriptions", IsAuthenticated(cache, logger, sc.Subscribe))
	v1.Delete("/subscriptions", IsAuthenticated(cache, logger, sc.Unsubscribe))
	v1.Get("/subscriptions", IsAuthenticated(cache, logger, sc.GetUserSubs))
	v1.Get("/subscriptions/posts", IsAuthenticated(cache, logger, sc.GetPostFromSub))

	r.Mount("/v1", v1)
	return r
}
