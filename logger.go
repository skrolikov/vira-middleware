package middleware

import (
	"context"
	"net/http"

	log "github.com/skrolikov/vira-logger"
)

func ContextLogger(base *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := base.WithContext(r.Context())
			ctx := context.WithValue(r.Context(), loggerKey, logger)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
