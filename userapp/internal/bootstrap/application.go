package bootstrap

import (
	"userapp/internal/application"
	"userapp/internal/application/api"
	"userapp/internal/application/service"
)

func RegisterService() {
	appContainer.RegisterService("templateService", new(service.Template))
	appContainer.RegisterService("activityLogService", new(service.ActivityLog))
	appContainer.RegisterService("authService", new(service.Auth))
	appContainer.RegisterService("authServiceLoggable", new(service.AuthLoggable))
	appContainer.RegisterService("pocketService", new(service.Pocket))
}

func RegisterApi() {
	appContainer.RegisterService("healthCheckHandler", new(api.HealthCheckHandler))
	appContainer.RegisterService("templateHandler", new(api.TemplateHandler))
	appContainer.RegisterService("authHandler", new(api.AuthHandler))
	appContainer.RegisterService("pocketHandler", new(api.PocketHandler))

	appContainer.RegisterService("api", new(application.Api))
}
