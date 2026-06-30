package querier

import (
	"context"

	"seedapp/internal/domain/seed"

	sq "github.com/Masterminds/squirrel"
)

type SeedQuerier interface {
	GetSchemas(ctx context.Context) (string, []interface{}, error)
	CreateSeedHistory(ctx context.Context, seed *seed.Seed) (string, []interface{}, error)
	UpdateSeedHistory(ctx context.Context, seed *seed.Seed) (string, []interface{}, error)
	GetLatestSeedVersionByType(ctx context.Context, schema string, seedType seed.SeedType) (string, []interface{}, error)
}

type Seed struct {
	SqlBuilder sq.StatementBuilderType `inject:"sqlBuilder"`
}

func (m *Seed) GetSchemas(ctx context.Context) (string, []interface{}, error) {
	query := m.SqlBuilder.Select("schema_name").
		From("information_schema.schemata").
		Where(sq.NotEq{"schema_name": []string{"information_schema", "pg_catalog", "pg_toast", "default", "public", "main"}})

	return query.ToSql()
}

func (m *Seed) CreateSeedHistory(ctx context.Context, seed *seed.Seed) (string, []interface{}, error) {
	query := m.SqlBuilder.Insert("main.seeds").
		Columns("schema", "version", "status", "error", "started_at", "finished_at").
		Values(seed.Schema, seed.Version, seed.Status, seed.Error, seed.StartedAt, seed.FinishedAt)

	query = query.Suffix("RETURNING id")

	return query.ToSql()
}

func (m *Seed) UpdateSeedHistory(ctx context.Context, seed *seed.Seed) (string, []interface{}, error) {
	query := m.SqlBuilder.Update("main.seeds").
		Set("version", seed.Version).
		Set("status", seed.Status).
		Set("entity_processed", seed.EntityProcessed).
		Set("error", seed.Error).
		Set("finished_at", seed.FinishedAt).
		Where("schema = ?", seed.Schema).
		Where("id = ?", seed.IDSerial)

	return query.ToSql()
}

// Get latest seed version on schema default
func (m *Seed) GetLatestSeedVersionByType(ctx context.Context, schema string, seedType seed.SeedType) (string, []interface{}, error) {
	query := m.SqlBuilder.Select("version", "entity_processed").
		From("main.seeds").
		Where(sq.Eq{"schema": schema,
			"status": seed.SEED_STATUS_SUCCESS.String()}).
		OrderBy("created_at DESC").Limit(1)

	return query.ToSql()
}
