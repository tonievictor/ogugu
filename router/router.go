package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/swaggo/http-swagger/v2"
)

func Routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	v1 := chi.NewRouter()
	v1.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	v1.Get("/swagger/*", httpSwagger.Handler())

	r.Mount("/v1", v1)
	return r
}
