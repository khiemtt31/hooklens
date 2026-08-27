package httpapi

import (
	"net/http"
	"server/internal/health"
)

type HealthHandler struct {
	service *health.Service
}

func NewHealthHandler(service *health.Service) *HealthHandler {
	return &HealthHandler{
		service: service,
	}
}

func (h *HealthHandler) Get(w http.ResponseWriter, r *http.Request) {
	report := h.service.Check()

	writeResponse(w, http.StatusOK, "health check passed", report)
}
