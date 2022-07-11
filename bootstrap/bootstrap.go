package bootstrap

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpServer "github.com/charmLd/token-generator-api/http/server"
	"github.com/charmLd/token-generator-api/util/config"
	"github.com/charmLd/token-generator-api/util/container"

	"log"
)

// Start starts the application bootstrap process
func Start() {

	// parse all configurations
	cfg := config.Parse()
	//timezone
	location, err := time.LoadLocation(cfg.AppConf.Timezone)
	if err != nil {
		log.Fatal(`Cannot load time location`, err, location)
	}

	// resolve the container using parsed configurations
	ctr := container.Resolve(cfg, location)

	//Initialize HTTP and HTTPs server
	httpSrv := httpServer.NewHTTPServer(cfg, ctr)
	errs := httpSrv.Init()

	// This will run forever until channel receives error
	select {
	case err := <-errs:
		log.Printf("Could not start serving service due to (error: %s)", err)
	}

	c := make(chan os.Signal, 1)

	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.

	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT)

	// Block until we receive our signal
	sig := <-c

	log.Println("bootstrap.init.Start", fmt.Sprintf("Received Signal: %s", sig))

	// Shutdown Http server
	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	httpSrv.ShutDown(ctx)

	<-ctx.Done()

	// Destruct other resources and stop the service
	Destruct(ctr)

	log.Println("bootstrap.init.Start", "Service shutted down gracefully...")
	os.Exit(0)

}

// Destruct Gracefully closes all additional resources.
//
// NOTE: For the container make sure to call destructors in the
// opposite order of initializing them
func Destruct(ctr *container.Container) {

	//todo Gracefully Shutdown the db connection

	// log.Println("bootstrap.Destruct", "Closing database connections...")
	// ctr.Adapters.MySQL.Destruct()

}
