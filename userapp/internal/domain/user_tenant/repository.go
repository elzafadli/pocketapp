package user_tenant

import "context"

//go:generate mockgen -destination=mocks/repository.go -package=mocks -source=repository.go
type Repository interface {
	Create(ctx context.Context, item *UserTenant) error
}
