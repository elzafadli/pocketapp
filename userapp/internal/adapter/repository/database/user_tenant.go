package database

import (
	"context"
	"fmt"

	"userapp/internal/adapter/repository"
	"userapp/internal/adapter/repository/database/querier"
	"userapp/internal/domain/user_tenant"

	"github.com/runsystemid/golog"
)

type UserTenantRepository struct {
	DB      repository.Sqlx            `inject:"database"`
	Querier querier.UserTenantQuerier  `inject:"userTenantQuerier"`
}

func (r *UserTenantRepository) Startup() error { return nil }
func (r *UserTenantRepository) Shutdown() error { return nil }

func (r *UserTenantRepository) Create(ctx context.Context, data *user_tenant.UserTenant) error {
	query, args, err := r.Querier.Create(data)
	if err != nil {
		golog.Error(ctx, fmt.Sprintf("Error build query create user_tenant: %s", err.Error()), err)
		return user_tenant.ErrDataCreate
	}

	_, err = r.DB.ExecContext(ctx, query, args...)
	if err != nil {
		golog.Error(ctx, fmt.Sprintf("Error create user_tenant: %s", err.Error()), err)
		return user_tenant.ErrDataCreate
	}

	return nil
}
