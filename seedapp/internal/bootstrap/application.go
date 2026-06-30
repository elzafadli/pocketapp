package bootstrap

import (
	"seedapp/internal/application"
	"seedapp/internal/application/api"
	"seedapp/internal/application/service"
)

func RegisterService() {
	appContainer.RegisterService("migrationService", new(service.Migration))
	appContainer.RegisterService("seedService", new(service.Seed))
}

func RegisterApi() {
	appContainer.RegisterService("healthCheckHandler", new(api.HealthCheckHandler))
	appContainer.RegisterService("migrationHandler", new(api.MigrationHandler))
	appContainer.RegisterService("seedHandler", new(api.SeedHandler))
	appContainer.RegisterService("api", new(application.Api))
}

