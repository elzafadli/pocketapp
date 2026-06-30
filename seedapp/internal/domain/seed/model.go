package seed

import (
	"seedapp/internal/domain/shared/entity"
	"time"

	"github.com/lib/pq"

	"gopkg.in/guregu/null.v4"
)

type Seed struct {
	entity.Entity
	Schema          string         `json:"schema" db:"schema"`
	Version         string         `json:"version" db:"version"`
	Status          string         `json:"status" db:"status"`
	Error           string         `json:"error" db:"error"`
	EntityProcessed pq.StringArray `json:"entity_processed" db:"entity_processed"`
	StartedAt       time.Time      `json:"started_at" db:"started_at"`
	FinishedAt      time.Time      `json:"finished_at" db:"finished_at"`
	SeedType        SeedType       `json:"seed_type" db:"seed_type"`
}

// for delete demo data
type SeedDefaultData struct {
	DocumentNumber        string      `json:"document_number" db:"document_number"`
	VoucherDocumentNumber null.String `json:"voucher_number" db:"voucher_document_number"`
}

type SeedingBudgetPlanningItems struct {
	DocumentNumber string `db:"document_number"`
	DetailNumber   string `db:"detail_number"`
}
