package router

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"

	"ogugu/controllers/common/response"
)

type Session struct {
	UserID     string
	SessionID  string
	CreatedAt  time.Time
	ExpiryTime time.Time
}

var tracer = otel.Tracer("middleware")

const AuthSession = "AuthSession"

func IsAuthenticated(cache *redis.Client, log *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			spanctx, span := tracer.Start(r.Context(), "Is authenticated middleware")
			defer span.End()

			auth := r.Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Bearer ") {
				log.Error("session id not found in authorization header")
				response.Error(w, "You are not logged in", http.StatusUnauthorized, log)
				return
			}

			token := strings.TrimPrefix(auth, "Bearer ")
			if token == "" {
				log.Error("session id not found in suffix to Bearer")
				response.Error(w, "You are not logged in", http.StatusUnauthorized, log)
				return
			}

			value, err := cache.Get(context.Background(), token).Result()
			if err != nil && err != redis.Nil {
				log.Error("provided session id does not exist in cache", zap.Error(err))
				response.Error(w, "You are not logged in", http.StatusUnauthorized, log)
				return
			}

			var session Session
			err = json.Unmarshal([]byte(value), &session)
			if err != nil {
				log.Error("session token cannot be validated into a session", zap.Error(err))
				response.Error(w, "Invalid token", http.StatusInternalServerError, log)
				return
			}

			if session.ExpiryTime.Before(time.Now()) {
				log.Error("provided session token has expired")
				response.Error(w, "You are not logged in", http.StatusUnauthorized, log)
				return
			}

			ctx := context.WithValue(spanctx, AuthSession, session)
			req := r.WithContext(ctx)

			next.ServeHTTP(w, req)
		}

		return http.HandlerFunc(fn)
	}
}
