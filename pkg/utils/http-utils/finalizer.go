package httputils

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
)

type Status struct {
	Status string `json:"status"`
}

type ErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Status       int           `json:"status"`
	Message      string        `json:"message"`
	Details      []ErrorDetail `json:"details"`
	LogTimestamp string        `json:"log_timestamp"`
	RequestID    string        `json:"request_id"`
}

func WriteResponse(w http.ResponseWriter, status int, message string, err error, data interface{}) interface{} {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err != nil {
		errorResponse := ErrorResponse{
			Status:       status,
			Message:      message,
			Details:      []ErrorDetail{{Field: "general", Message: err.Error()}},
			LogTimestamp: time.Now().Format(time.RFC3339),
			RequestID:    uuid.New().String(),
		}

		w.WriteHeader(status)
		v, _ := jsoniter.Marshal(errorResponse)
		w.Write(v)

		return errorResponse
	}

	if data != nil {
		w.WriteHeader(status)
		v, _ := jsoniter.Marshal(data)
		w.Write(v)
		return data
	}

	statusResponse := Status{
		Status: message,
	}

	w.WriteHeader(http.StatusOK)
	v, _ := jsoniter.Marshal(statusResponse)
	w.Write(v)

	return statusResponse
}
