package user

import "context"

//go:generate mockgen -destination=mocks/repository.go -package=mocks -source=repository.go
type Repository interface {
	Create(ctx context.Context, item *User) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, item *User) error
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context, filter map[string]interface{}) ([]*User, error)
	Count(ctx context.Context, filter map[string]interface{}) (uint64, error)
}
