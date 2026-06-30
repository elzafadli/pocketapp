package application

import (
	"seedapp/internal/adapter/rest"
	"seedapp/internal/application/api"
	"seedapp/internal/pkg/custommiddleware"
)

type Api struct {
	App                *rest.Fiber                            `inject:"fiber"`
	AuthMiddleware     custommiddleware.AuthMiddlewareService `inject:"authMiddleware"`
	HealthCheckHandler api.HealthCheckAPI                     `inject:"healthCheckHandler"`
	MigrationHandler   api.MigrationAPI                       `inject:"migrationHandler"`
	SeedHandler        api.SeedAPI                            `inject:"seedHandler"`
}

func (a *Api) Startup() error {
	a.App.Get("/ping", a.HealthCheckHandler.Ping)
	a.App.Get("/ready", a.HealthCheckHandler.Ready)
	a.App.Get("/version", a.HealthCheckHandler.Version)

	v1 := a.App.Group("/v1")

	v1.Post("/migrate", a.AuthMiddleware.BasicAuthProtection(), a.MigrationHandler.RunMigration)
	v1.Post("/seed", a.AuthMiddleware.BasicAuthProtection(), a.SeedHandler.RunSeed)

	return nil
}

func (a *Api) Shutdown() error {
	return nil
}
