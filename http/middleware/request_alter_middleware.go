package middleware

import (
	"context"
	"github.com/charmLd/token-generator-api/domain/globals"
	"net/http"

	"github.com/google/uuid"
)

// RequestAlterMiddleware alters the request
type RequestAlterMiddleware struct{}

// NewRequestAlterMiddleware creates a new instance of RequestAlterMiddleware
func NewRequestAlterMiddleware() *RequestAlterMiddleware {
	return &RequestAlterMiddleware{}
}

// Middleware executes middleware rules of RequestAlterMiddleware
func (rtm *RequestAlterMiddleware) Middleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Get uuid from request header, if not exits create a new UUID version 4
		contextUUID := r.Header.Get("x-correlation-id")
		if contextUUID == "" {
			contextUUID = uuid.New().String()
		}

		ctx := context.WithValue(r.Context(), globals.UUIDKey, contextUUID)

		r = r.WithContext(ctx)
		// call the next handler, which can be another middleware in the chain, or the final handler
		next.ServeHTTP(w, r)
	})
}
