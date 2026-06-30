package bootstrap

import (
	"userapp/internal/adapter/agentapp"
	"userapp/internal/adapter/monitor"
	"userapp/internal/adapter/kaisel"
	"userapp/internal/adapter/repository"
	"userapp/internal/adapter/repository/cache"
	"userapp/internal/adapter/repository/database"
	"userapp/internal/adapter/repository/database/querier"
	"userapp/internal/adapter/rest"

	sq "github.com/Masterminds/squirrel"
)

func RegisterAgent() {
	appContainer.RegisterService("agentClient", new(agentapp.Client))
}

func RegisterMonitor() {
	appContainer.RegisterService("monitorClient", new(monitor.Monitor))
}

func RegisterDatabase() {
	appContainer.RegisterService("database", new(repository.SqlxDB))
}

func RegisterCache() {
	appContainer.RegisterService("cache", new(repository.Cache))
}

func RegisterRest() {
	appContainer.RegisterService("fiber", new(rest.Fiber))
}

func RegisterKaisel() {
	appContainer.RegisterService("kaisel", new(kaisel.Kaisel))
}

func RegisterSQLBuilder() {
	// Squirrel with PostgreSQL dollar-sign placeholders
	sqlBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	appContainer.RegisterService("sqlBuilder", sqlBuilder)
}

func RegisterQuerier() {
	appContainer.RegisterService("activityLogQuerier", new(querier.ActivityLog))
	appContainer.RegisterService("userQuerier", new(querier.User))
	appContainer.RegisterService("userTenantQuerier", new(querier.UserTenant))
	appContainer.RegisterService("pocketQuerier", new(querier.Pocket))
}

func RegisterRepository() {
	appContainer.RegisterService("templateRepository", new(database.TemplateRepository))
	appContainer.RegisterService("templateCacheRepository", new(cache.TemplateRepository))
	appContainer.RegisterService("activityLogRepository", new(database.ActivityLogRepository))
	appContainer.RegisterService("userRepository", new(database.UserRepository))
	appContainer.RegisterService("userTenantRepository", new(database.UserTenantRepository))
	appContainer.RegisterService("pocketRepository", new(database.PocketRepository))
}
