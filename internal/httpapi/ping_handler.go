package httpapi

import (
	"encoding/json"
	"net/http"

	"server/internal/ping"
)

type PingHandler struct {
	service *ping.Service
}

func NewPingHandler(service *ping.Service) *PingHandler {
	return &PingHandler{
		service: service,
	}
}

func (h *PingHandler) Post(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var request ping.Request

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		GlobalErrorHandler(w, http.StatusBadRequest, "invalid request body")

		return
	}

	if request.Pong == "" {
		GlobalErrorHandler(w, http.StatusBadRequest, "pong is required")

		return
	}

	response := h.service.Ping(request.Pong)

	writeResponse(w, http.StatusOK, "pong received", response)
}
