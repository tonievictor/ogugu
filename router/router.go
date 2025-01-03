package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

func Routes(logger *zap.Logger) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	v1 := chi.NewRouter()

	v1.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte("Up and ready to rumble!!!\n"))
	})

	v1.Get("/swagger/*", httpSwagger.Handler())

	r.Mount("/v1", v1)
	return r
}
