package sqlx

import (
	"context"
	"database/sql"
	"errors"
	"seedapp/internal/adapter/repository"
	"seedapp/internal/adapter/repository/sqlx/querier"
	"seedapp/internal/domain/tenant"

	"github.com/runsystemid/golog"
)

type TenantRepository struct {
	DB      repository.Sqlx       `inject:"database"`
	Querier querier.TenantQuerier `inject:"tenantQuerier"`
}

func (r *TenantRepository) GetAll(ctx context.Context, filter map[string]any) ([]*tenant.Tenant, error) {
	query, args, err := r.Querier.GetAll(ctx, filter)
	if err != nil {
		golog.Error(ctx, "error querier get all tenant", err)
		return nil, NewErrQuery(err)
	}

	var res []*tenant.Tenant
	if err = r.DB.SelectContext(ctx, &res, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, tenant.ErrTenantNotFound
		}
		golog.Error(ctx, "error exec getting schemas", err)
		return nil, NewErrQuery(err)
	}

	return res, nil
}

func (r *TenantRepository) Get(ctx context.Context, tenantCode string) (*tenant.Tenant, error) {
	query, args, err := r.Querier.Get(ctx, tenantCode)
	if err != nil {
		golog.Error(ctx, "error querier get tenant", err)
		return nil, NewErrQuery(err)
	}

	var res tenant.Tenant
	if err = r.DB.GetContext(ctx, &res, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, tenant.ErrTenantNotFound
		}
		golog.Error(ctx, "error exec getting tenant", err)
		return nil, NewErrQuery(err)
	}

	return &res, nil
}

func (r *TenantRepository) Create(ctx context.Context, tenant *tenant.Tenant) error {
	query, args, err := r.Querier.Create(ctx, tenant)
	if err != nil {
		golog.Error(ctx, "error querier create tenant", err)
		return NewErrQuery(err)
	}

	_, err = r.DB.ExecContext(ctx, query, args...)
	if err != nil {
		golog.Error(ctx, "error exec create tenant", err)
		return NewErrQuery(err)
	}

	return nil
}
