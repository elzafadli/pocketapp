package sqlx

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"seedapp/internal/adapter/repository"
	"seedapp/internal/adapter/repository/sqlx/querier"
	"seedapp/internal/domain/migration"
	"seedapp/scripts/migrations"

	migrate "github.com/rubenv/sql-migrate"
	"github.com/runsystemid/golog"
)

type MigrationRepository struct {
	DB      repository.Sqlx          `inject:"database"`
	Querier querier.MigrationQuerier `inject:"migrationQuerier"`
}

func (r *MigrationRepository) GetSchemas(ctx context.Context) ([]string, error) {
	query, args, err := r.Querier.GetSchemas(ctx)
	if err != nil {
		golog.Error(ctx, "error querier get schemas", err)
		return nil, NewErrQuery(err)
	}

	var res []string
	if err = r.DB.SelectContext(ctx, &res, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, migration.ErrSchemaNotFound
		}
		golog.Error(ctx, "error exec getting schemas", err)
		return nil, NewErrQuery(err)
	}

	return res, nil
}

func (r *MigrationRepository) CreateMigrationHistory(ctx context.Context, migration *migration.Migration) error {
	query, args, err := r.Querier.CreateMigrationHistory(ctx, migration)
	if err != nil {
		golog.Error(ctx, "error querier create migration history", err)
		return NewErrQuery(err)
	}

	_, err = r.DB.ExecContext(ctx, query, args...)
	if err != nil {
		golog.Error(ctx, "error exec create migration history", err)
		return NewErrQuery(err)
	}

	return nil
}

func (r *MigrationRepository) UpdateMigrationHistory(ctx context.Context, migration *migration.Migration) error {
	query, args, err := r.Querier.UpdateMigrationHistory(ctx, migration)
	if err != nil {
		golog.Error(ctx, "error querier update migration history", err)
		return NewErrQuery(err)
	}

	_, err = r.DB.ExecContext(ctx, query, args...)
	if err != nil {
		golog.Error(ctx, "error exec update migration history", err)
		return NewErrQuery(err)
	}

	return nil
}

func (r *MigrationRepository) RunMigrations(ctx context.Context, migrationType migration.MigrationType, schema string) (totalApplied int, err error) {
	trimSchema := strings.TrimSpace(schema)
	trimSchema = strings.ReplaceAll(trimSchema, "\"", "")

	if migrationType == migration.MIGRATION_TYPE_TENANT {
		_, err = r.DB.ExecContext(ctx, fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS "%s"`, trimSchema))
		if err != nil {
			golog.Error(ctx, fmt.Sprintf("%s: error creating schema", trimSchema), err)
			return 0, NewErrQuery(err)
		}
	}

	// Set search path to the target schema
	_, err = r.DB.ExecContext(ctx, fmt.Sprintf(`SET search_path TO "%s"`, trimSchema))
	if err != nil {
		golog.Error(ctx, "error setting search path", err)
		return 0, NewErrQuery(err)
	}

	migs := []*migrate.Migration{}
	if migrationType == migration.MIGRATION_TYPE_MAIN {
		migs = migrations.MainMigrations()
		migrate.SetSchema(migration.SCHEMA_MAIN)
	} else if migrationType == migration.MIGRATION_TYPE_TENANT {
		migs = migrations.TenantMigrations(schema)
		migrate.SetSchema(fmt.Sprintf(`"%s"`, trimSchema))
	}

	migSource := &migrate.MemoryMigrationSource{Migrations: migs}

	res, err := migrate.Exec(r.DB.GetDB().DB, repository.POSTGRES_DRIVER, migSource, migrate.Up)
	if err != nil {
		// Check if error is an "already exists" error that we can safely ignore
		if r.DB.SqlxDBIsObjectMigrationAlreadyExists(err) {
			golog.Warn(ctx, fmt.Sprintf("%s: migration skipped (object already exists): %v", schema, err))
		} else {
			// For other errors, proceed with normal error handling
			golog.Error(ctx, fmt.Sprintf("%s: error exec migrations", schema), err)
			if res > 0 {
				golog.Info(ctx, fmt.Sprintf("%s: rolling back %d migrations", schema, res))
				rollRes, _err := migrate.ExecMax(r.DB.GetDB().DB, repository.POSTGRES_DRIVER, migSource, migrate.Down, res)
				if _err != nil {
					golog.Error(ctx, fmt.Sprintf("%s: error exec rollback migrations", schema), _err)
					return 0, NewErrQuery(_err)
				}
				golog.Info(ctx, fmt.Sprintf("%s: rolled back %d migrations", schema, rollRes))
				return 0, NewErrQuery(err)
			} else {
				return 0, NewErrQuery(err)
			}
		}
	}

	if res > 0 {
		golog.Info(ctx, fmt.Sprintf("%s: applied %d migrations", schema, res))
	} else {
		golog.Info(ctx, fmt.Sprintf("%s: no migrations applied", schema))
	}

	return res, nil
}

func (r *MigrationRepository) GetLatestMigrationVersion(ctx context.Context, schema string) (string, error) {
	query, args, err := r.Querier.GetLatestMigrationVersion(ctx, schema)
	if err != nil {
		golog.Error(ctx, "error querier get latest migration version", err)
		return "", NewErrQuery(err)
	}

	var res string
	if err = r.DB.QueryRowxContext(ctx, query, args...).Scan(&res); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", migration.ErrMigrationNotFound
		}
		golog.Error(ctx, "error exec get latest migration version", err)
		return "", NewErrQuery(err)
	}

	return res, nil
}

func (r *MigrationRepository) GetLatestMigrationStatus(ctx context.Context, schema string) (string, error) {
	query, args, err := r.Querier.GetLatestMigrationStatus(ctx, schema)
	if err != nil {
		golog.Error(ctx, "error querier get latest migration status", err)
		return "", NewErrQuery(err)
	}

	var res string
	if err = r.DB.QueryRowxContext(ctx, query, args...).Scan(&res); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", migration.ErrMigrationNotFound
		}
		golog.Error(ctx, "error exec get latest migration status", err)
		return "", NewErrQuery(err)
	}

	return res, nil
}
