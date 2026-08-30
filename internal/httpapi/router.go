package httpapi

import "net/http"

func NewRouter(
	healthHandler *HealthHandler,
	pingHandler *PingHandler,
	inboxHandler *InboxHandler,
	eventHandler *EventHandler,
) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/health", methodHandler(http.MethodGet, healthHandler.Get))
	mux.HandleFunc("/api/inboxes", inboxMethodHandler(inboxHandler))
	mux.HandleFunc("/api/inboxes/{id}/events", eventMethodHandler(eventHandler))
	mux.HandleFunc("/ping", methodHandler(http.MethodPost, pingHandler.Post))
	mux.HandleFunc("/", notFoundHandler)

	return mux
}

func eventMethodHandler(handler *EventHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.Get(w, r)
		case http.MethodPost:
			handler.Post(w, r)
		default:
			w.Header().Set("Allow", "GET, POST")
			GlobalErrorHandler(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func inboxMethodHandler(handler *InboxHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.Get(w, r)
		case http.MethodPost:
			handler.Post(w, r)
		default:
			w.Header().Set("Allow", "GET, POST")
			GlobalErrorHandler(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}
