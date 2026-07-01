package pocket

import (
	"context"

	"userapp/internal/domain/shared/identity"
)

//go:generate mockgen -destination=mocks/repository.go -package=mocks -source=repository.go
type Repository interface {
	Create(ctx context.Context, schema string, data *PocketItem) (*PocketItem, error)
	Update(ctx context.Context, schema string, data *PocketItem) error
	GetByID(ctx context.Context, schema string, id identity.ID) (*PocketItem, error)
	Delete(ctx context.Context, schema string, id identity.ID) error
	List(ctx context.Context, schema string, filter map[string]interface{}) ([]*PocketItem, error)
	Count(ctx context.Context, schema string, filter map[string]interface{}) (uint64, error)
	GetSummary(ctx context.Context, schema string) (*PocketSummary, error)
}
