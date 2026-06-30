package migrations

import (
	migrationmain "seedapp/scripts/migrations/main"
	migrationtenant "seedapp/scripts/migrations/tenant"

	migrate "github.com/rubenv/sql-migrate"
)

// sql-migration will reorder by migration id
//
// Make sure new migration id is greater than the last migration id
func MainMigrations() []*migrate.Migration {
	return []*migrate.Migration{
		migrationmain.CreateTableUsers(),
		migrationmain.CreateUserTenant(),
		migrationmain.CreateTenant(),
	}
}

func TenantMigrations(schema string) []*migrate.Migration {
	return []*migrate.Migration{
		migrationtenant.CreatePocketItems(schema),
	}
}
