package bootstrap

import (
	"monitorapp/internal/adapter/repository"
	"monitorapp/internal/adapter/repository/cache"
	"monitorapp/internal/adapter/repository/database"
	"monitorapp/internal/adapter/repository/database/querier"
	"monitorapp/internal/adapter/rest"

	sq "github.com/Masterminds/squirrel"
)

func RegisterDatabase() {
	appContainer.RegisterService("database", new(repository.SqlxDB))
}

func RegisterCache() {
	appContainer.RegisterService("cache", new(repository.Cache))
}

func RegisterRest() {
	appContainer.RegisterService("fiber", new(rest.Fiber))
}

func RegisterSQLBuilder() {
	// Squirrel with PostgreSQL dollar-sign placeholders
	sqlBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	appContainer.RegisterService("sqlBuilder", sqlBuilder)
}

func RegisterQuerier() {
	appContainer.RegisterService("activityLogQuerier", new(querier.ActivityLog))
}

func RegisterRepository() {
	appContainer.RegisterService("templateRepository", new(database.TemplateRepository))
	appContainer.RegisterService("templateCacheRepository", new(cache.TemplateRepository))
	appContainer.RegisterService("activityLogRepository", new(database.ActivityLogRepository))
}
