package application

import (
	"monitorapp/internal/adapter/rest"
	"monitorapp/internal/application/api"
	"monitorapp/internal/pkg/custommiddleware"
)

type Api struct {
	App                *rest.Fiber                            `inject:"fiber"`
	HealthCheckHandler api.HealthCheckAPI                     `inject:"healthCheckHandler"`
	TemplateHandler    api.TemplateAPI                        `inject:"templateHandler"`
	ActivityLogHandler api.ActivityLogAPI                     `inject:"activityLogHandler"`
	AuthMiddleware     custommiddleware.AuthMiddlewareService `inject:"authMiddleware"`
}

func (a *Api) Startup() error {
	a.App.Get("/ping", a.HealthCheckHandler.Ping)
	a.App.Get("/ready", a.HealthCheckHandler.Ready)
	a.App.Get("/version", a.HealthCheckHandler.Version)

	v1 := a.App.Group("/v1")
	v1.Get("/template", a.AuthMiddleware.BasicAuthProtection(), a.TemplateHandler.Get)
	v1.Post("/template", a.AuthMiddleware.BasicAuthProtection(), a.TemplateHandler.Create)
	v1.Put("/template/:id", a.AuthMiddleware.BasicAuthProtection(), a.TemplateHandler.Update)
	v1.Get("/template/:id", a.AuthMiddleware.BasicAuthProtection(), a.TemplateHandler.Detail)

	api.ActivityLogRoutes(v1, a.AuthMiddleware.BasicAuthProtection(), a.ActivityLogHandler)

	return nil
}

func (a *Api) Shutdown() error {
	return nil
}
