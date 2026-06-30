package sqlx

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"

	"seedapp/internal/adapter/repository"
	"seedapp/internal/adapter/repository/sqlx/querier"
	"seedapp/internal/domain/seed"
	"seedapp/scripts/seeds"

	"github.com/runsystemid/golog"
)

type SeedRepository struct {
	DB      repository.Sqlx     `inject:"database"`
	Querier querier.SeedQuerier `inject:"seedQuerier"`
	Seeder  seeds.SeederService `inject:"seeder"`
}

func (r *SeedRepository) GetSchemas(ctx context.Context) ([]string, error) {
	query, args, err := r.Querier.GetSchemas(ctx)
	if err != nil {
		golog.Error(ctx, "error querier get schemas", err)
		return nil, NewErrQuery(err)
	}

	var res []string
	if err = r.DB.SelectContext(ctx, &res, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, seed.ErrSchemaNotFound
		}
		golog.Error(ctx, "error exec getting schemas", err)
		return nil, NewErrQuery(err)
	}

	return res, nil
}

func (r *SeedRepository) CreateSeedHistory(ctx context.Context, seed *seed.Seed) (int64, error) {
	query, args, err := r.Querier.CreateSeedHistory(ctx, seed)
	if err != nil {
		golog.Error(ctx, "error querier create seed history", err)
		return 0, NewErrQuery(err)
	}

	var id int64
	if err = r.DB.GetContext(ctx, &id, query, args...); err != nil {
		golog.Error(ctx, "error exec create seed history", err)
		return 0, NewErrQuery(err)
	}

	return id, nil
}

func (r *SeedRepository) UpdateSeedHistory(ctx context.Context, seed *seed.Seed) error {
	query, args, err := r.Querier.UpdateSeedHistory(ctx, seed)
	if err != nil {
		golog.Error(ctx, "error querier update seed history", err)
		return NewErrQuery(err)
	}

	_, err = r.DB.ExecContext(ctx, query, args...)
	if err != nil {
		golog.Error(ctx, "error exec update seed history", err)
		return NewErrQuery(err)
	}

	return nil
}

func (r *SeedRepository) RunSeeds(ctx context.Context, tenantType seed.SeedTenantType, schema string, entityAlreadyProcessed map[string]bool) ([]string, error) {
	var entityProcessed []string
	trimSchema := strings.TrimSpace(schema)
	trimSchema = strings.ReplaceAll(trimSchema, "\"", "")

	// Set search path to the target schema
	_, err := r.DB.ExecContext(ctx, fmt.Sprintf(`SET search_path TO "%s"`, trimSchema))
	if err != nil {
		golog.Error(ctx, "error setting search path", err)
		return nil, NewErrQuery(err)
	}

	allSeeds := r.Seeder.AllSeeds()

	sortedSeedName := r.Seeder.AllSeedsName(allSeeds)

	for _, seedName := range sortedSeedName {
		seedFunc := allSeeds[seedName]

		processed := entityAlreadyProcessed[seedName]
		if !processed {
			err = seedFunc(ctx, tenantType, trimSchema, r.DB.GetDB())
			if err != nil {
				golog.Error(ctx, fmt.Sprintf("%s: error exec seeds", seedName), err)
				return entityProcessed, NewErrQuery(err)
			}
		}
		entityProcessed = append(entityProcessed, seedName)
	}

	golog.Info(ctx, fmt.Sprintf("%s: applied seeds", schema))

	sort.Strings(entityProcessed)
	return entityProcessed, nil
}

func (r *SeedRepository) GetLatestSeedVersionByType(ctx context.Context, schema string, seedType seed.SeedType) (*seed.Seed, error) {
	query, args, err := r.Querier.GetLatestSeedVersionByType(ctx, schema, seedType)
	if err != nil {
		golog.Error(ctx, "error querier get seed version", err)
		return nil, NewErrQuery(err)
	}

	var res seed.Seed
	if err = r.DB.QueryRowxContext(ctx, query, args...).StructScan(&res); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, seed.ErrSeedNotFound
		}
		golog.Error(ctx, "error exec get seed version seeding", err)
		return nil, NewErrQuery(err)
	}

	return &res, nil
}
