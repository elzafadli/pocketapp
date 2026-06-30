package menu

import (
	"errors"
	"time"
)

type Menu struct {
	ID        string     `db:"id" json:"id"`
	Name      string     `db:"name" json:"name"`
	Path      string     `db:"path" json:"path"`
	Icon      string     `db:"icon" json:"icon"`
	ParentID  *string    `db:"parent_id" json:"parent_id"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
}

type CreateMenuRequest struct {
	Name     string  `json:"name" validate:"required,max=255"`
	Path     string  `json:"path" validate:"max=255"`
	Icon     string  `json:"icon" validate:"max=255"`
	ParentID *string `json:"parent_id"`
}

type UpdateMenuRequest struct {
	ID       string  `json:"id" validate:"required"`
	Name     string  `json:"name" validate:"max=255"`
	Path     string  `json:"path" validate:"max=255"`
	Icon     string  `json:"icon" validate:"max=255"`
	ParentID *string `json:"parent_id"`
}

var (
	ErrDataNotFound      = errors.New("menu data not found")
	ErrDataAlreadyExists = errors.New("menu data already exists")
	ErrDataCreate        = errors.New("failed to create menu")
	ErrDataUpdate        = errors.New("failed to update menu")
	ErrDataDelete        = errors.New("failed to delete menu")
	ErrDataGet           = errors.New("failed to get menu")
)
