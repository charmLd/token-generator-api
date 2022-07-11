package router

import (
	"github.com/charmLd/token-generator-api/http/controllers"
	"github.com/charmLd/token-generator-api/http/middleware"
	"github.com/charmLd/token-generator-api/http/transport/response"
	"github.com/charmLd/token-generator-api/util/container"
	"github.com/gorilla/mux"
	"net/http"
)

// Init initializes the router.
func InitHTTP(container *container.Container) *mux.Router {

	// create new router
	r := mux.NewRouter()

	// initialize middleware
	requestCheckerMiddleware := middleware.NewRequestCheckerMiddleware(container)
	requestAlterMiddleware := middleware.NewRequestAlterMiddleware()
	authenticationMiddleware := middleware.NewAuthenticationMiddleware(container)
	r.Use(authenticationMiddleware.Middleware)
	r.Use(requestCheckerMiddleware.Middleware)
	r.Use(requestAlterMiddleware.Middleware)

	// initialize controllers
	baseController := controllers.NewBaseController(container)

	// api info
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		response.Send(
			w,
			nil,
			http.StatusOK)
	}).Methods(http.MethodGet)

	//r.HandleFunc("/user/create", baseController.CreateUser).Methods(http.MethodPost)

	//since we need to authenticate the admin endpoints we generate jwt to token from this login endpoint
	r.HandleFunc("/auth/login", baseController.EmailLoginControllerFunc).Methods(http.MethodPost)

	//-- admin endpoints
	// generate token

	r.HandleFunc("/token/generate", baseController.AdminTokenGenerateControllerFun).Methods(http.MethodPost)

	//disable token
	r.HandleFunc("/token/revoke", baseController.AdminTokenRevokeControllerFun).Methods(http.MethodPost)

	//retrieve disable and active token
	r.HandleFunc("/user/{user_id}/token/fetch", baseController.AdminTokenFetchControllerFunc).Methods(http.MethodGet)

	return r
}

// Init initializes the router.
func InitHTTPs(container *container.Container) *mux.Router {

	// create new router
	r := mux.NewRouter()

	// initialize middleware
	requestCheckerMiddleware := middleware.NewRequestCheckerMiddleware(container)
	throttlecheckMiddleware := middleware.NewThrottleMiddleware(container)
	requestAlterMiddleware := middleware.NewRequestAlterMiddleware()

	r.Use(requestCheckerMiddleware.Middleware)
	r.Use(throttlecheckMiddleware.Middleware)
	r.Use(requestAlterMiddleware.Middleware)

	// initialize controllers
	baseController := controllers.NewBaseController(container)

	// api info
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		response.Send(
			w,
			nil,
			http.StatusOK)
	}).Methods(http.MethodGet)

	//validate invite token  // public api - no need to authenticate
	r.HandleFunc("/auth/validate", baseController.PublicTokenValidateCOntrollerFunc).Methods(http.MethodPost)
	return r
}
