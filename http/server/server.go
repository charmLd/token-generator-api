package server

import (
	"context"
	//"fmt"
	"github.com/charmLd/token-generator-api/http/router"
	"github.com/charmLd/token-generator-api/util/config"
	"github.com/charmLd/token-generator-api/util/container"
	"log"
	"net/http"
	"strconv"
	"time"
)

// HTTPServer implementes the base type for Http server
type HTTPServer struct {
	serverHTTP *http.Server

	config    *config.Config
	container *container.Container
}

// NewHTTPServer creates a new HttpServer
func NewHTTPServer(config *config.Config, container *container.Container) *HTTPServer {
	return &HTTPServer{
		serverHTTP: nil,
		config:     config,
		container:  container,
	}
}

// Init initializes the server
func (srv *HTTPServer) Init() chan error {
	errs := make(chan error)
	// initialize the router
	r := router.InitHTTP(srv.container)

	rhttps := router.InitHTTPs(srv.container)

	internalAddress := srv.config.AppConf.Host + ":" + strconv.Itoa(srv.config.AppConf.InternalPort)

	serverHTTP := &http.Server{
		Addr: internalAddress,

		// good practice to set timeouts to avoid Slowloris attacks
		WriteTimeout: time.Second * 60,
		ReadTimeout:  time.Second * 60,
		IdleTimeout:  time.Second * 60,

		// pass our instance of gorilla/mux in
		Handler: r,
	}
	publicAddress := srv.config.AppConf.Host + ":" + strconv.Itoa(srv.config.AppConf.PublicPort)
	serverHTTPs := &http.Server{
		Addr: publicAddress,

		// good practice to set timeouts to avoid Slowloris attacks
		WriteTimeout: time.Second * 60,
		ReadTimeout:  time.Second * 60,
		IdleTimeout:  time.Second * 60,

		// pass our instance of gorilla/mux in
		Handler: rhttps,
	}
	// run our server in a goroutine so that it doesn't block
	go func() {
		log.Printf("Staring HTTP service on %s ", internalAddress)
		err := serverHTTP.ListenAndServe()
		if err != nil {
			log.Println("http.server.Init", err)
			errs <- err
		}
	}()

	// run our server in a goroutine so that it doesn't block HTTPs

	go func() {
		log.Printf("Staring HTTPS service on %s ", publicAddress)
		if err := serverHTTPs.ListenAndServeTLS("./config/test.crt", "./config/test.key"); err != nil {
			errs <- err
		}
	}()

	srv.serverHTTP = serverHTTP
	//log.Println("http.server.Init", fmt.Sprintf("HTTP server listening on %s", publicAddress))

	return nil

}

//todo check this
// ShutDown releases all http connections gracefully and shut down the server
func (srv *HTTPServer) ShutDown(ctx context.Context) {

	go func() {

		log.Println("http.server.ShutDown", "Stopping HTTP Server")
		srv.serverHTTP.SetKeepAlivesEnabled(false)

		// Doesn't block if no connections, but will otherwise wait
		// until the timeout deadline.
		err := srv.serverHTTP.Shutdown(ctx)
		if err != nil {
			log.Println("http.server.ShutDown", "Unable to stop HTTP server")

		}
	}()
}
