package middleware

import (
	"gitlab.com/otqee/otqee-be/internal/logger"
	"go.uber.org/zap"
	"net/http"
)

// AccessLogMiddleware is a middleware for access log
func AccessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		logger.Logger.Debug("access",
			zap.String("from", r.RemoteAddr),
			zap.String("url", r.URL.Path),
		)
	})
}
