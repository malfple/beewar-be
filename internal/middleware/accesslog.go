package middleware

import (
	"fmt"
	"net/http"
)

// AccessLogMiddleware is a middleware for access log
func AccessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		fmt.Println("accessed")
	})
}
