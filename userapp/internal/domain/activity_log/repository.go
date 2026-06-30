package activity_log

import "context"

//go:generate mockgen -destination=mocks/repository.go -package=mocks -source=repository.go
type Repository interface {
	Create(ctx context.Context, item *ActivityLog) (*ActivityLog, error)
	GetByID(ctx context.Context, id string) (*ActivityLog, error)
	Update(ctx context.Context, item *ActivityLog) error
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context, filter map[string]interface{}) ([]*ActivityLog, error)
	Count(ctx context.Context, filter map[string]interface{}) (uint64, error)
}
