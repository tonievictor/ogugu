package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"

	"ogugu/response"
)

type Session struct {
	UserID     string
	SessionID  string
	CreatedAt  time.Time
	ExpiryTime time.Time
}

var tracer = otel.Tracer("middleware")

const AuthSession = "AuthSession"

func IsAuthenticated(next http.HandlerFunc, cache *redis.Client, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		spanctx, span := tracer.Start(r.Context(), "Is authenticated middleware")
		defer span.End()

		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			response.Error(w, "You are not logged in", http.StatusUnauthorized, errors.New("session id not found"), log)
			return
		}

		token := strings.TrimPrefix(auth, "Bearer ")
		if token == "" {
			response.Error(w, "You are not logged in", http.StatusUnauthorized, errors.New("session id not found"), log)
			return
		}

		value, err := cache.Get(context.Background(), token).Result()
		if err == redis.Nil {
			response.Error(w, "You are not logged in", http.StatusUnauthorized, errors.New("invalid session token"), log)
			return
		}

		var session Session
		err = json.Unmarshal([]byte(value), &session)
		if err != nil {
			response.Error(w, "Invalid token", http.StatusInternalServerError, errors.New("Could not validate auth token"), log)
			return
		}

		if session.ExpiryTime.Before(time.Now()) {
			response.Error(w, "You are not logged in", http.StatusUnauthorized, errors.New("invalid session token"), log)
			return
		}

		ctx := context.WithValue(spanctx, AuthSession, session)
		req := r.WithContext(ctx)

		next(w, req)
	}
}
