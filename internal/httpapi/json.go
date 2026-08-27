package httpapi

import (
	"encoding/json"
	"log"
	"net/http"
)

type APIResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Data       any    `json:"data,omitempty"`
}

func writeResponse(w http.ResponseWriter, statusCode int, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := APIResponse{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("failed to encode JSON response: %v", err)
	}
}

func GlobalErrorHandler(w http.ResponseWriter, statusCode int, message string) {
	writeResponse(w, statusCode, message, nil)
}
