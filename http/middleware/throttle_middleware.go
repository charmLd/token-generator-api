package middleware

import (
	"fmt"
	"net/http"

	e "github.com/charmLd/token-generator-api/http/error"
	"github.com/charmLd/token-generator-api/http/error/types"
	"github.com/charmLd/token-generator-api/util/container"
	log "github.com/sirupsen/logrus"

	"strings"
	"time"
)

// RequestAlterMiddleware alters the request
type ThrottleMiddleware struct {
	Container *container.Container
}

// NewRequestAlterMiddleware creates a new instance of RequestAlterMiddleware
func NewThrottleMiddleware(ctr *container.Container) *ThrottleMiddleware {
	return &ThrottleMiddleware{
		Container: ctr,
	}
}

// Middleware executes middleware rules of RequestAlterMiddleware
func (tm *ThrottleMiddleware) Middleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

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

		JwtPayloadData, err := tm.Container.Adapters.Token.ValidateLoginJWToken(r.Context(), token)
		if err != nil {
			log.Error(r.Context(), "access header Bearer JWT token error : ", err)
			e.Handle(r.Context(), (&types.UnAuthorizeError{}).New("access header Bearer JWT token error"), w)
			return
		}
		lastlogintime, err := tm.Container.Repositories.UserRepository.GetLastLoginTime(r.Context(), fmt.Sprint(JwtPayloadData.UserID))
		if err != nil {
			log.Error(r.Context(), "error : ", err)
			e.Handle(r.Context(), (&types.UnAuthorizeError{}).New("unauthorized activity"), w)
			return
		}

		fmt.Println("diff time : ", time.Since(lastlogintime))

		if time.Since(lastlogintime) < time.Duration(tm.Container.ThrottledTime)*time.Minute { //it will not allow the user to do valid requests until the throttled time period if over
			//check the time tifference between last login times and
			//if it exceed return error
			e.Handle(r.Context(), (&types.UnAuthorizeError{}).New("exceed the request rate for the client | Throttled"), w)
		}

		next.ServeHTTP(w, r)
	})
}
