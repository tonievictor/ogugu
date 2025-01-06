package response

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type Response struct {
	message string      `json:"message"`
	data    interface{} `json: "data,omitempty"`
}

func Error(w http.ResponseWriter, message string, status int, data interface{}, log *zap.Logger) {
	log.Error("RESPONSE", zap.String("message", message))
	res := Response{
		message: message,
		data:    data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Error("RESPONSE", zap.String("Error sending response", err.Error()))
	}
}

func Success(w http.ResponseWriter, message string, status int, data interface{}, log *zap.Logger) {
	log.Info("RESPONSE", zap.String("message", message))
	res := Response{
		message: message,
		data:    data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Error("RESPONSE", zap.String("Error sending response", err.Error()))
	}
}
