package bootstrap

import (
	"monitorapp/internal/application"
	"monitorapp/internal/application/api"
	"monitorapp/internal/application/grpc"
	"monitorapp/internal/application/service"
)

func RegisterService() {
	appContainer.RegisterService("templateService", new(service.Template))
	appContainer.RegisterService("activityLogService", new(service.ActivityLog))
}

func RegisterApi() {
	appContainer.RegisterService("healthCheckHandler", new(api.HealthCheckHandler))
	appContainer.RegisterService("templateHandler", new(api.TemplateHandler))
	appContainer.RegisterService("activityLogHandler", new(api.ActivityLogHandler))
	appContainer.RegisterService("activityLogGrpcHandler", new(grpc.ActivityLogHandler))

	appContainer.RegisterService("api", new(application.Api))
}
