package response

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"ogugu/internal/models"
)

func Error(w http.ResponseWriter, message string, status int, log *zap.Logger) {
	res := Response{
		Message: message,
		Data:    "",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Error("RESPONSE", zap.String("Error sending response", err.Error()))
	}
}

func Success(w http.ResponseWriter, message string, status int, data any, log *zap.Logger) {
	res := Response{
		Message: message,
		Data:    data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Error("RESPONSE", zap.String("Error sending response", err.Error()))
	}
}

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type RssFeed struct {
	Message string
	Data    models.RssFeed
}

type Post struct {
	Message string
	Data    models.Post
}

type Posts struct {
	Message string
	Data    []models.Post
}

type User struct {
	Message string
	Data    models.User
}

type UserWithAuth struct {
	Message string
	Data    models.UserWithAuth
}

type RssFeeds struct {
	Message string
	Data    []models.RssFeed
}

type Subscription struct {
	Message string
	Data    models.Subscription
}

type FeedPosts struct {
	Message string
	Data    []models.Post
}
