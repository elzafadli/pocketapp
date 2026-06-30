package querier

import (
	"context"

	"seedapp/internal/domain/migration"

	sq "github.com/Masterminds/squirrel"
)

type MigrationQuerier interface {
	GetSchemas(ctx context.Context) (string, []interface{}, error)
	CreateMigrationHistory(ctx context.Context, migration *migration.Migration) (string, []interface{}, error)
	UpdateMigrationHistory(ctx context.Context, migration *migration.Migration) (string, []interface{}, error)
	GetLatestMigrationVersion(ctx context.Context, schema string) (string, []interface{}, error)
	GetLatestMigrationStatus(ctx context.Context, schema string) (string, []interface{}, error)
}

type Migration struct {
	SqlBuilder sq.StatementBuilderType `inject:"sqlBuilder"`
}

func (m *Migration) GetSchemas(ctx context.Context) (string, []interface{}, error) {
	query := m.SqlBuilder.Select("schema_name").
		From("information_schema.schemata").
		Where(sq.NotEq{"schema_name": []string{"information_schema", "pg_catalog", "pg_toast", "default", "public", "main"}})

	return query.ToSql()
}

func (m *Migration) CreateMigrationHistory(ctx context.Context, migration *migration.Migration) (string, []interface{}, error) {
	query := m.SqlBuilder.Insert("main.migrations").
		Columns("schema", "version", "status", "error", "started_at", "finished_at").
		Values(migration.Schema, migration.Version, migration.Status, migration.Error, migration.StartedAt, migration.FinishedAt)

	return query.ToSql()
}

func (m *Migration) UpdateMigrationHistory(ctx context.Context, migration *migration.Migration) (string, []interface{}, error) {
	query := m.SqlBuilder.Update("main.migrations").
		Set("status", migration.Status).
		Set("error", migration.Error).
		Set("finished_at", migration.FinishedAt).
		Where("schema = ? AND version = ?", migration.Schema, migration.Version)

	return query.ToSql()
}

// Get latest migration version on schema default
func (m *Migration) GetLatestMigrationVersion(ctx context.Context, schema string) (string, []interface{}, error) {
	query := m.SqlBuilder.Select("version").
		From("main.migrations").
		Where(sq.Eq{"schema": schema, "status": "success"}).
		OrderBy("created_at DESC").Limit(1)

	return query.ToSql()
}

func (m *Migration) GetLatestMigrationStatus(ctx context.Context, schema string) (string, []interface{}, error) {
	query := m.SqlBuilder.Select("status").
		From("main.migrations").
		Where(sq.Eq{"schema": schema}).
		OrderBy("created_at DESC").Limit(1)

	return query.ToSql()
}
