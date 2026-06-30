package application

import (
	"userapp/internal/adapter/rest"
	"userapp/internal/application/api"
	"userapp/internal/pkg/custommiddleware"
)

type Api struct {
	App                *rest.Fiber                            `inject:"fiber"`
	HealthCheckHandler api.HealthCheckAPI                     `inject:"healthCheckHandler"`
	TemplateHandler    api.TemplateAPI                        `inject:"templateHandler"`
	AuthHandler        api.AuthAPI                            `inject:"authHandler"`
	PocketHandler      api.PocketAPI                          `inject:"pocketHandler"`
	AuthMiddleware     custommiddleware.AuthMiddlewareService `inject:"authMiddleware"`
}

func (a *Api) Startup() error {
	a.App.Get("/ping", a.HealthCheckHandler.Ping)
	a.App.Get("/ready", a.HealthCheckHandler.Ready)
	a.App.Get("/version", a.HealthCheckHandler.Version)

	apiGroup := a.App.Group("/api")
	authGroup := apiGroup.Group("/auth")

	api.AuthRoutes(authGroup, a.AuthMiddleware, a.AuthHandler)
	api.PocketRoutes(apiGroup, a.AuthMiddleware, a.PocketHandler)

	return nil
}

func (a *Api) Shutdown() error {
	return nil
}
