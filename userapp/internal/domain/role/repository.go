package role

import "context"

//go:generate mockgen -destination=mocks/repository.go -package=mocks -source=repository.go
type Repository interface {
	Create(ctx context.Context, item *Role) (*Role, error)
	GetByID(ctx context.Context, id string) (*Role, error)
	Update(ctx context.Context, item *Role) error
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context, filter map[string]interface{}) ([]*Role, error)
	Count(ctx context.Context, filter map[string]interface{}) (uint64, error)
}
