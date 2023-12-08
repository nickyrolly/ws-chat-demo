package middleware

import (
	"net/http"
)

func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Exercise 4.1
		// please complete this block to add verifier for user token
	})
}
