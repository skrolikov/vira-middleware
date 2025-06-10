package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type ctxKeyRequestID struct{}

func RequestID() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := uuid.NewString()
			ctx := context.WithValue(r.Context(), ctxKeyRequestID{}, id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
