package middleware

import (
	"log"
	"net/http"
	"strconv"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userHeader := r.Header.Get("userId")
		if userHeader == "" {
			log.Printf("userID Header not set")
			http.Error(w, `Unauthorized (401)`, http.StatusUnauthorized)
			return
		}

		_, err := strconv.Atoi(userHeader)
		if err != nil {
			log.Printf("Invalid `userId` Header: %v", userHeader)
			http.Error(w, "Internal Server Error", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}
