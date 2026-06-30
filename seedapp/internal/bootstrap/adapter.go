package bootstrap

import (
	"seedapp/internal/adapter/cache"
	"seedapp/internal/adapter/repository"
	"seedapp/internal/adapter/repository/sqlx"
	"seedapp/internal/adapter/repository/sqlx/querier"
	"seedapp/internal/adapter/rest"
	"seedapp/scripts/seeds"

	sq "github.com/Masterminds/squirrel"
)

func RegisterDatabase() {
	appContainer.RegisterService("database", new(repository.SqlxDB))
	appContainer.RegisterService("seeder", new(seeds.Seeder))
}

func RegisterRest() {
	appContainer.RegisterService("fiber", new(rest.Fiber))
}

func RegisterRepository() {
	appContainer.RegisterService("sqlBuilder", sq.StatementBuilder.PlaceholderFormat(sq.Dollar))

	appContainer.RegisterService("migrationQuerier", new(querier.Migration))
	appContainer.RegisterService("migrationRepository", new(sqlx.MigrationRepository))

	appContainer.RegisterService("seedQuerier", new(querier.Seed))
	appContainer.RegisterService("seedRepository", new(sqlx.SeedRepository))

	appContainer.RegisterService("tenantQuerier", new(querier.Tenant))
	appContainer.RegisterService("tenantRepository", new(sqlx.TenantRepository))
}






func RegisterCache() {
	appContainer.RegisterService("cache", new(cache.Redis))
}
