package querier

import (
	"context"
	"seedapp/internal/domain/tenant"

	sq "github.com/Masterminds/squirrel"
)

var (
	TenantTable = "main.tenants"
)

type TenantQuerier interface {
	GetAll(ctx context.Context, filter map[string]any) (string, []interface{}, error)
	Get(ctx context.Context, tenantCode string) (string, []interface{}, error)
	Create(ctx context.Context, tenant *tenant.Tenant) (string, []interface{}, error)
}

type Tenant struct {
	SqlBuilder sq.StatementBuilderType `inject:"sqlBuilder"`
}

func (m *Tenant) GetAll(ctx context.Context, filter map[string]any) (string, []interface{}, error) {
	query := m.SqlBuilder.Select("tenant_code", "tenant_name").
		From(TenantTable + " t")

	if migratedOnly, ok := filter["migrated_only"].(bool); ok && migratedOnly {
		query = query.Join(PG_NAMESPACE + " AS n ON n.nspname = t.tenant_code::text")
	}

	query = query.OrderBy("created_at DESC")

	return query.ToSql()
}

func (m *Tenant) Get(ctx context.Context, tenantCode string) (string, []interface{}, error) {
	query := m.SqlBuilder.Select("tenant_code", "tenant_name").
		From(TenantTable + " t").
		Where(sq.Eq{"t.tenant_code": tenantCode})

	return query.ToSql()
}

func (m *Tenant) Create(ctx context.Context, tenant *tenant.Tenant) (string, []interface{}, error) {
	return m.SqlBuilder.Insert(TenantTable).
		Columns("tenant_code", "tenant_name").
		Values(tenant.TenantCode, tenant.TenantName).
		ToSql()
}
