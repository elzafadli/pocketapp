package menu

import "context"

//go:generate mockgen -destination=mocks/repository.go -package=mocks -source=repository.go
type Repository interface {
	Create(ctx context.Context, item *Menu) (*Menu, error)
	GetByID(ctx context.Context, id string) (*Menu, error)
	Update(ctx context.Context, item *Menu) error
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context, filter map[string]interface{}) ([]*Menu, error)
	Count(ctx context.Context, filter map[string]interface{}) (uint64, error)
}
