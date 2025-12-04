package api

import (
	"net/http"
	"os"
)

func APIKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected_key := os.Getenv("MICHO_API_KEY")
		if expected_key == "" {
			http.Error(w, "Server configuration error: API Key not set", http.StatusInternalServerError)
			return
		}
		client_key := r.Header.Get("X-Api-Key")
		if client_key != expected_key {
			http.Error(w, "Invalid API Key", http.StatusUnauthorized)
		}
		next.ServeHTTP(w, r)
	})
}
