package seed

import (
	"context"
)

//go:generate mockgen -destination=mocks/repository.go -package=mocks -source=repository.go
type Repository interface {
	GetSchemas(ctx context.Context) ([]string, error)
	CreateSeedHistory(ctx context.Context, seed *Seed) (int64, error)
	UpdateSeedHistory(ctx context.Context, seed *Seed) error
	RunSeeds(ctx context.Context, tenantType SeedTenantType, schema string, entityAlreadyProcessed map[string]bool) ([]string, error)
	GetLatestSeedVersionByType(ctx context.Context, schema string, seedType SeedType) (*Seed, error)
}
