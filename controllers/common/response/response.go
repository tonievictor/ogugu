package response

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"ogugu/models"
)

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type RssFeed struct {
	Message string
	Data    models.RssFeed
}

type RssFeeds struct {
	Message string
	Data    []models.RssFeed
}

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
