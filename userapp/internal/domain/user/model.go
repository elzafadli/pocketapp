package user

import (
	"errors"
	"time"
)

type User struct {
	ID              string     `db:"id" json:"id"`
	Name            string     `db:"name" json:"name"`
	TenantDefault   string     `db:"tenant_default" json:"tenant_default"`
	ActiveIndicator string     `db:"active_indicator" json:"active_indicator"`
	Email           string     `db:"email" json:"email"`
	Password        string     `db:"password" json:"-"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt       *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,max=255"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=6"`
}

type UpdateUserRequest struct {
	ID       string `json:"id" validate:"required"`
	Name     string `json:"name" validate:"max=255"`
	Email    string `json:"email" validate:"omitempty,email,max=255"`
	Password string `json:"password" validate:"omitempty,min=6"`
}

var (
	ErrDataNotFound      = errors.New("user data not found")
	ErrDataAlreadyExists = errors.New("user data already exists")
	ErrDataCreate        = errors.New("failed to create user")
	ErrDataUpdate        = errors.New("failed to update user")
	ErrDataDelete        = errors.New("failed to delete user")
	ErrDataGet           = errors.New("failed to get user")
)
