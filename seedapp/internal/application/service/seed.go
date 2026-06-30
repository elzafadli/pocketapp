package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"seedapp/config"
	"seedapp/internal/application/model"
	"seedapp/internal/domain/seed"
	"seedapp/internal/domain/shared/entity"
	"seedapp/internal/domain/tenant"
	"seedapp/scripts/seeds"

	"github.com/runsystemid/golog"
)

//go:generate mockgen -destination=mocks/seed.go -package=mocks -source=seed.go
type SeedService interface {
	RunSeeds(ctx context.Context, request *model.SeedRequest) (*model.SeedResponse, error)
}

type Seed struct {
	Seeder           seeds.SeederService     `inject:"seeder"`
	SeedRepo         seed.Repository         `inject:"seedRepository"`
	Config           *config.Config          `inject:"config"`
	TenantRepository tenant.TenantRepository `inject:"tenantRepository"`
}

func (s *Seed) RunSeeds(ctx context.Context, request *model.SeedRequest) (*model.SeedResponse, error) {
	schemas, err := s.getTargetSchemas(ctx, request.Schemas)
	if err != nil {
		return nil, err
	}

	response := &model.SeedResponse{
		Success: make([]string, len(schemas)),
		Failed:  make([]model.SeedError, len(schemas)),
	}

	for _, sch := range schemas {
		golog.Info(ctx, fmt.Sprintf("running seeds for schema %s", sch))
		if err := s.processSingleSchema(ctx, request.TenantType, sch, response); err != nil {
			return nil, err
		}
	}

	return response, nil
}

func (s *Seed) getTargetSchemas(ctx context.Context, requestSchemas []string) ([]string, error) {
	if len(requestSchemas) > 0 {
		return requestSchemas, nil
	}

	return s.SeedRepo.GetSchemas(ctx)
}

func (s *Seed) processSingleSchema(ctx context.Context, tenantType seed.SeedTenantType, schema string, response *model.SeedResponse) error {
	_seed, err := s.SeedRepo.GetLatestSeedVersionByType(ctx, schema, seed.SEED_TYPE_SEEDING)
	if err != nil && !errors.Is(err, seed.ErrSeedNotFound) {
		return err
	}

	allSeedsName := s.Seeder.AllSeedsName(s.Seeder.AllSeeds())
	latestVersion := s.Seeder.HashSeed(allSeedsName)
	if _seed != nil && latestVersion == _seed.Version {
		return nil
	}

	seedRecord := &seed.Seed{
		Entity:          entity.NewEntity(),
		Schema:          schema,
		Version:         latestVersion,
		EntityProcessed: make([]string, 0),
		Status:          seed.SEED_STATUS_RUNNING.String(),
		StartedAt:       time.Now(),
		SeedType:        seed.SEED_TYPE_SEEDING,
	}

	idSeedHistory, err := s.SeedRepo.CreateSeedHistory(ctx, seedRecord)
	if err != nil {
		return err
	}

	seedRecord.IDSerial = idSeedHistory

	entityAlreadyProcessed := make(map[string]bool)
	if _seed != nil {
		for _, seedName := range _seed.EntityProcessed {
			entityAlreadyProcessed[seedName] = true
		}
	}

	entityProcessed, err := s.SeedRepo.RunSeeds(ctx, tenantType, schema, entityAlreadyProcessed)
	if err != nil {
		seedRecord.Status = seed.SEED_STATUS_FAILED.String()
		seedRecord.Version = s.Seeder.HashSeed(entityProcessed)
		seedRecord.EntityProcessed = entityProcessed
		seedRecord.Error = err.Error()
		response.Failed = append(response.Failed, model.SeedError{
			Schema: schema,
			Error:  err.Error(),
		})
	} else {
		seedRecord.Status = seed.SEED_STATUS_SUCCESS.String()
		seedRecord.EntityProcessed = entityProcessed
		response.Success = append(response.Success, schema)
	}
	seedRecord.FinishedAt = time.Now()

	if err := s.SeedRepo.UpdateSeedHistory(ctx, seedRecord); err != nil {
		return err
	}

	return nil
}
