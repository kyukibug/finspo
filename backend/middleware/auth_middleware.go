package middleware

import (
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("userId") == "" {
			http.Error(w, "User ID is not set", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}
