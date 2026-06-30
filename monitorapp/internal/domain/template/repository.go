package template

import (
	"context"

	"monitorapp/internal/domain/shared/identity"
)

//go:generate mockgen -destination=mocks/repository.go -package=mocks -source=repository.go
type Repository interface {
	Create(ctx context.Context, data *Template) (*Template, error)
	Update(ctx context.Context, data *Template) error
	GetByID(ctx context.Context, id identity.ID) (*Template, error)
}
