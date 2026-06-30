package bootstrap

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"monitorapp/config"
	"monitorapp/internal/adapter/rest"
	"monitorapp/internal/pkg/custommiddleware"
	"monitorapp/internal/pkg/validator"

	"github.com/runsystemid/golog"
	"github.com/runsystemid/gontainer"
)

var appContainer = gontainer.New()

func Run(conf *config.Config) {
	appContainer.RegisterService("config", conf)

	// Initialize struct validator
	appContainer.RegisterService("validator", validator.NewGoValidator())
	appContainer.RegisterService("authMiddleware", new(custommiddleware.AuthMiddleware))

	bootstrapContext, cancel := context.WithCancel(context.Background())
	golog.Info(bootstrapContext, "Serving...")

	// Register adapter
	RegisterDatabase()
	RegisterCache()
	RegisterRest()
	RegisterSQLBuilder()
	RegisterQuerier()
	RegisterRepository()

	// Register application
	RegisterService()
	RegisterApi()

	// Startup the container
	if err := appContainer.Ready(); err != nil {
		golog.Panic(bootstrapContext, "Failed to populate service", err)
	}

	// Start server
	fiberApp := appContainer.GetServiceOrNil("fiber").(*rest.Fiber)
	errs := make(chan error, 3)
	go func() {
		golog.Info(bootstrapContext, fmt.Sprintf("Listening on port :%d", conf.Http.Port))
		errs <- fiberApp.Listen(fmt.Sprintf(":%d", conf.Http.Port))
	}()

	golog.Info(bootstrapContext, "Your app started")
	printLogo(conf.Http.Port)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errs:
		if err != nil {
			golog.Error(bootstrapContext, "Server error", err)
		}
		cancel()
	case sig := <-quit:
		golog.Info(bootstrapContext, fmt.Sprintf("Signal termination received: %v", sig))
		cancel()
	}

	<-bootstrapContext.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	golog.Info(shutdownCtx, "Shutting down server...")

	// Gracefully shutdown Fiber server first
	if err := fiberApp.Shutdown(); err != nil {
		golog.Error(shutdownCtx, "Error shutting down Fiber server", err)
	}

	golog.Info(shutdownCtx, "Cleaning up resources...")
	appContainer.Shutdown()

	golog.Info(shutdownCtx, "Bye")
}

func printLogo(port int) {
	fmt.Printf(`
██████╗░░█████╗░██╗░░██╗██╗░░░██╗
██╔══██╗██╔══██╗██║░██╔╝██║░░░██║
██████╔╝███████║█████═╝░██║░░░██║
██╔═══╝░██╔══██║██╔═██╗░██║░░░██║
██║░░░░░██║░░██║██║░╚██╗╚██████╔╝
╚═╝░░░░░╚═╝░░╚═╝╚═╝░░╚═╝░╚═════╝░
Running on port %d
`+"\n", port)
}
