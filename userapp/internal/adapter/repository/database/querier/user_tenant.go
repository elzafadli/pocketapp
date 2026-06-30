package querier

import (
	"userapp/internal/domain/user_tenant"

	sq "github.com/Masterminds/squirrel"
)

const (
	USER_TENANT_TABLE = "main.user_tenants"
)

type UserTenantQuerier interface {
	Create(data *user_tenant.UserTenant) (string, []interface{}, error)
}

type UserTenant struct {
	SQLBuilder sq.StatementBuilderType `inject:"sqlBuilder"`
}

func (q *UserTenant) Create(data *user_tenant.UserTenant) (string, []interface{}, error) {
	return q.SQLBuilder.
		Insert(USER_TENANT_TABLE).
		Columns("user_code", "tenant_code", "active_indicator", "created_by", "created_at", "updated_by", "updated_at").
		Values(data.UserCode, data.TenantCode, data.ActiveIndicator, data.CreatedBy, data.CreatedAt, data.UpdatedBy, data.UpdatedAt).
		ToSql()
}
