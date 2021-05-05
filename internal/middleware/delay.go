package middleware

import (
	"net/http"
	"time"
)

// DelayMiddleware is used to simulate latency
func DelayMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second)
		next.ServeHTTP(w, r)
	})
}
