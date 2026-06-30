package role_menu

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Actions []string

func (a Actions) Value() (driver.Value, error) {
	if a == nil {
		return "[]", nil
	}
	return json.Marshal(a)
}

func (a *Actions) Scan(value interface{}) error {
	if value == nil {
		*a = []string{}
		return nil
	}
	var b []byte
	switch v := value.(type) {
	case string:
		b = []byte(v)
	case []byte:
		b = v
	default:
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, a)
}

type RoleMenu struct {
	RoleID    string     `db:"role_id" json:"role_id"`
	MenuID    string     `db:"menu_id" json:"menu_id"`
	Actions   Actions    `db:"actions" json:"actions"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
}

type CreateRoleMenuRequest struct {
	RoleID  string   `json:"role_id" validate:"required"`
	MenuID  string   `json:"menu_id" validate:"required"`
	Actions []string `json:"actions"`
}

type UpdateRoleMenuRequest struct {
	RoleID  string   `json:"role_id" validate:"required"`
	MenuID  string   `json:"menu_id" validate:"required"`
	Actions []string `json:"actions"`
}

var (
	ErrDataNotFound      = errors.New("role menu data not found")
	ErrDataAlreadyExists = errors.New("role menu data already exists")
	ErrDataCreate        = errors.New("failed to create role menu")
	ErrDataUpdate        = errors.New("failed to update role menu")
	ErrDataDelete        = errors.New("failed to delete role menu")
	ErrDataGet           = errors.New("failed to get role menu")
)
