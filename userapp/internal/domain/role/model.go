package role

import (
	"errors"
	"time"

	"userapp/internal/domain/menu"
)

type RoleMenuDetail struct {
	menu.Menu
	Actions []string `json:"actions"`
}

type Role struct {
	ID          string            `db:"id" json:"id"`
	Name        string            `db:"name" json:"name"`
	Description string            `db:"description" json:"description"`
	CreatedAt   time.Time         `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time         `db:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time        `db:"deleted_at" json:"deleted_at"`
	Menus       []*RoleMenuDetail `json:"menus"`
}

type CreateRoleRequest struct {
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description"`
}

type UpdateRoleRequest struct {
	ID          string `json:"id" validate:"required"`
	Name        string `json:"name" validate:"max=255"`
	Description string `json:"description"`
}

var (
	ErrDataNotFound      = errors.New("role data not found")
	ErrDataAlreadyExists = errors.New("role data already exists")
	ErrDataCreate        = errors.New("failed to create role")
	ErrDataUpdate        = errors.New("failed to update role")
	ErrDataDelete        = errors.New("failed to delete role")
	ErrDataGet           = errors.New("failed to get role")
)
