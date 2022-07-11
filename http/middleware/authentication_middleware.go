package middleware

import (
	"fmt"
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
		fmt.Println("toekn:", token)

		JwtPayloadData, err := a.Container.Adapters.Token.DecodeAuthToken(r.Context(), token)
		fmt.Println(JwtPayloadData, "  ", err)
		//JwtPayloadData, err := DecodeJwtPayload(r.Context(), token)
		//fmt.Println("JWT Payload Data :", JwtPayloadData)

		if err != nil {
			log.Error(r.Context(), "access header Bearer JWT token error : ", err)
			e.Handle(r.Context(), (&types.UnAuthorizeError{}).New("access header Bearer JWT token error"), w)
			return
		}
		if JwtPayloadData.UserRole != "admin" {
			log.Error(r.Context(), "unauthorized user  ", err)
			e.Handle(r.Context(), (&types.UnAuthorizeError{}).New("unauthorized user"), w)
			return
		}
		/*
		   check if the token has necessary roles to do the admin functionalities

		   //role - admin can call all endpoints - here only cosider admin role
		*/

		//validate the authorization token to authenticate the user or client
		/* isValid, err := a.Auth.ValidateAuthToken(r.Context(), token)
		if err != nil {
			log.Error(r.Context(), "JWT token validate error : ", err)
			e.Handle(r.Context(), (&types.UnAuthorizeError{}).New("JWT token validate error"), w)
			return

		}

		if !isValid {
			log.Error(r.Context(), "", "JWT token unauthorized : ")
			e.Handle(r.Context(), (&types.UnAuthorizeError{}).New("JWT token unauthorized : "), w)
			return
		} */
		//check if the client has the necessary permissions to do the admin functionalities unless return authentication error
		fmt.Println(JwtPayloadData)
		next.ServeHTTP(w, r)
	})
}

/* func DecodeJwtPayload(ctx context.Context, jwt string) (payload entities.JWTClaims, err error) {
	jwtSplit := strings.Split(jwt, ".")
	if len(jwtSplit) != 3 {
		log.Error(ctx, "Jwt Token Format Error :", jwt, err)
		return payload, err
	}
	data, err := base64.StdEncoding.WithPadding(base64.NoPadding).DecodeString(jwtSplit[1])
	if err != nil {
		log.Error(ctx, "Jwt base 64 decode error :", string(data), err)
		return payload, err
	}
	err = json.Unmarshal(data, &payload)
	if err != nil {
		log.Error(ctx, "Unmarshal Auth Token Error :", data, err)
		return payload, err
	}
	return payload, nil
}
*/
