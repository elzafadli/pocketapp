package migration

import (
	"seedapp/internal/domain/shared/entity"
	"time"
)

type Migration struct {
	entity.Entity
	Schema     string    `json:"schema" db:"schema"`
	Version    string    `json:"version" db:"version"`
	Status     string    `json:"status" db:"status"`
	Error      string    `json:"error" db:"error"`
	StartedAt  time.Time `json:"started_at" db:"started_at"`
	FinishedAt time.Time `json:"finished_at" db:"finished_at"`
}
