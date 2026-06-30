package role_menu

import "context"

//go:generate mockgen -destination=mocks/repository.go -package=mocks -source=repository.go
type Repository interface {
	Create(ctx context.Context, item *RoleMenu) (*RoleMenu, error)
	Get(ctx context.Context, roleID, menuID string) (*RoleMenu, error)
	Update(ctx context.Context, item *RoleMenu) error
	Delete(ctx context.Context, roleID, menuID string) error
	GetAll(ctx context.Context, filter map[string]interface{}) ([]*RoleMenu, error)
	Count(ctx context.Context, filter map[string]interface{}) (uint64, error)
}
