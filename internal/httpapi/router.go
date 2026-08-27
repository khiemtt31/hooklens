package httpapi

import "net/http"

func NewRouter(healthHandler *HealthHandler, pingHandler *PingHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/health", methodHandler(http.MethodGet, healthHandler.Get))
	mux.HandleFunc("/ping", methodHandler(http.MethodPost, pingHandler.Post))
	mux.HandleFunc("/", notFoundHandler)

	return mux
}
