package httpapi

import "net/http"

func methodHandler(method string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.Header().Set("Allow", method)
			GlobalErrorHandler(w, http.StatusMethodNotAllowed, "method not allowed")

			return
		}

		next(w, r)
	}
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	GlobalErrorHandler(w, http.StatusNotFound, "route not found")
}
