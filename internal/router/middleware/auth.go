package middleware

import (
	"net/http"
	"strings"

	"github.com/nickyrolly/ws-chat-demo/pkg/auth"
)

func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the JWT token from the request header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// The Authorization header is typically in the format "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate the JWT token
		err := auth.VerifyToken(tokenString)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// If the token is valid, pass the request to the next handler
		next.ServeHTTP(w, r)
	})
}
