package bootstrap

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"seedapp/config"
	"seedapp/internal/adapter/rest"
	"seedapp/internal/application/service"
	"seedapp/internal/pkg/custommiddleware"
	"seedapp/internal/pkg/validator"

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
	RegisterRest()
	RegisterRepository()
	RegisterCache()

	// Register application
	RegisterService()
	RegisterApi()

	// Startup the container
	if err := appContainer.Ready(); err != nil {
		golog.Panic(bootstrapContext, "Failed to populate service", err)
	}

	// Start server
	fiberApp := appContainer.GetServiceOrNil("fiber").(*rest.Fiber)
	errs := make(chan error, 2)
	go func() {
		golog.Info(bootstrapContext, fmt.Sprintf("Listening on port :%d", conf.Http.Port))
		errs <- fiberApp.Listen(fmt.Sprintf(":%d", conf.Http.Port))
	}()

	golog.Info(bootstrapContext, "Your app started")

	// Run migrations
	migrationService := appContainer.GetServiceOrNil("migrationService").(service.MigrationService)
	err := migrationService.RunDefaultMigration(bootstrapContext)
	if err != nil {
		golog.Panic(bootstrapContext, "Failed to run migrations", err)
	}

	errMigrate := migrationService.RunTenantsMigration(bootstrapContext)
	if errMigrate != nil {
		golog.Error(bootstrapContext, "Failed to run tenants migrations", errMigrate)
	}

	printLogo(conf.Http.Port)
	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		golog.Info(bootstrapContext, "Signal termination received")
		cancel()
	}()

	<-bootstrapContext.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	golog.Info(shutdownCtx, "Cleaning up resources...")

	appContainer.Shutdown()

	golog.Info(shutdownCtx, "Bye")
}

func printLogo(port int) {
	fmt.Printf(`Running on port %d`, port)
}
