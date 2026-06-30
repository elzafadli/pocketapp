// Package entity represents a shared kernel for domain model
package entity

import (
	"time"

	"seedapp/internal/domain/shared/identity"

	"gopkg.in/guregu/null.v4"
)

// Entity represents domain Entity
type Entity struct {
	ID        identity.ID `json:"id" db:"id"`
	IDSerial  int64       `json:"id_int" db:"-"`
	CreatedAt time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt time.Time   `json:"updated_at" db:"updated_at"`
	DeletedAt null.Time   `json:"deleted_at" db:"deleted_at"`
}

func NewEntity() Entity {
	now := time.Now()
	return Entity{
		ID:        identity.NewID(),
		CreatedAt: now,
		UpdatedAt: now,
	}
}
