package middleware

import (
	"net/http"
	"strings"

	e "github.com/charmLd/token-generator-api/http/error"
	"github.com/charmLd/token-generator-api/http/error/types"
	"github.com/charmLd/token-generator-api/util/container"
	log "github.com/sirupsen/logrus"
)

// AuthenticationMiddleware alters the request.
type AuthenticationMiddleware struct {
	Container *container.Container
}

// NewAuthenticationMiddleware creates a new instance of MetricMiddleware
func NewAuthenticationMiddleware(ctr *container.Container) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{
		Container: ctr,
	}
}

// The loggingResponseWriter is created embedding http.ResponseWriter
// https://golang.org/doc/effective_go.html#embedding
// https://ndersson.me/post/capturing_status_code_in_net_http/
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// Middleware executes middleware rules of AuthenticationMiddleware.
func (a *AuthenticationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		lrw := newLoggingResponseWriter(w)

		if (r.URL.Path == "/auth/login") || (r.URL.Path == "/") {
			next.ServeHTTP(lrw, r)
			return
		}
		///
		token := r.Header.Get("Authorization")
		if token == "" {
			e.Handle(r.Context(), (&types.ForbiddenError{}).New("authentication header invalid"), w)
			return
		}
		splitToken := strings.Split(token, "Bearer")
		if len(splitToken) != 2 {
			e.Handle(r.Context(), (&types.ForbiddenError{}).New("authentication header token not inserted"), w)
			return
		}
		token = strings.TrimSpace(splitToken[1])

		JwtPayloadData, err := a.Container.Adapters.Token.DecodeAuthToken(r.Context(), token)

		if err != nil {
			log.Error(r.Context(), "access header Bearer JWT token error : ", err)
			e.Handle(r.Context(), (&types.UnAuthorizeError{}).New("access header Bearer JWT token error"), w)
			return
		}
		if JwtPayloadData.UserRole != "admin" { //Only admin users can do admin functionalities
			log.Error(r.Context(), "unauthorized user  ", err)
			e.Handle(r.Context(), (&types.UnAuthorizeError{}).New("unauthorized user"), w)
			return
		}

		next.ServeHTTP(w, r)
	})
}
