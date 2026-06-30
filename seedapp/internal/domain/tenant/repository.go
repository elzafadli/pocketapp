package tenant

import (
	"context"
)

//go:generate mockgen -destination=mocks/repository.go -package=mocks -source=repository.go
type TenantRepository interface {
	GetAll(ctx context.Context, filter map[string]any) ([]*Tenant, error)
	Get(ctx context.Context, tenantCode string) (*Tenant, error)
	Create(ctx context.Context, tenant *Tenant) error
}
