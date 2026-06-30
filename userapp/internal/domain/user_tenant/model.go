package user_tenant

import (
	"errors"
	"time"
)

type UserTenant struct {
	UserCode        string    `db:"user_code"`
	TenantCode      string    `db:"tenant_code"`
	ActiveIndicator string    `db:"active_indicator"`
	CreatedBy       string    `db:"created_by"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedBy       string    `db:"updated_by"`
	UpdatedAt       time.Time `db:"updated_at"`
}

var (
	ErrDataCreate = errors.New("failed to create user_tenant")
)
