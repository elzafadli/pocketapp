package migration

import (
	"context"
)

//go:generate mockgen -destination=mocks/repository.go -package=mocks -source=repository.go
type Repository interface {
	GetSchemas(ctx context.Context) ([]string, error)
	CreateMigrationHistory(ctx context.Context, migration *Migration) error
	UpdateMigrationHistory(ctx context.Context, migration *Migration) error
	RunMigrations(ctx context.Context, migrationType MigrationType, schema string) (int, error)
	GetLatestMigrationVersion(ctx context.Context, schema string) (string, error)
	GetLatestMigrationStatus(ctx context.Context, schema string) (string, error)
}
