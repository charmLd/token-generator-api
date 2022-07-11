package middleware

import (
	"fmt"
	"net/http"

	"github.com/charmLd/token-generator-api/http/error"
	"github.com/charmLd/token-generator-api/http/error/types"
	"github.com/charmLd/token-generator-api/util/container"
)

// RequestCheckerMiddleware validates the request header
type RequestCheckerMiddleware struct {
	container     *container.Container
	omittedRoutes []string
}

// NewRequestCheckerMiddleware creates a new instance of RequestCheckerMiddleware
func NewRequestCheckerMiddleware(ctr *container.Container) *RequestCheckerMiddleware {

	omittedRoutes := []string{
		"/favicon.ico",
	}

	return &RequestCheckerMiddleware{
		container:     ctr,
		omittedRoutes: omittedRoutes,
	}
}

// Middleware executes middleware rules of RequestCheckerMiddleware
func (rtm *RequestCheckerMiddleware) Middleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		requestURI := r.RequestURI
		contentType := r.Header.Get("Content-Type")
		httpMethod := r.Method

		// skip omitted routes
		for _, v := range rtm.omittedRoutes {

			if v == requestURI {
				next.ServeHTTP(w, r)
				return
			}
		}

		// check content type
		if contentType != "application/json" && httpMethod != http.MethodDelete && httpMethod != http.MethodGet {

			err := types.MiddlewareError{}
			error.Handle(r.Context(), err.New(fmt.Sprintf("API only accepts JSON, '%s' is given", contentType), 100, ""), w)

			return
		}

		// call the next handler, which can be another middleware in the chain, or the final handler
		next.ServeHTTP(w, r)
	})
}
