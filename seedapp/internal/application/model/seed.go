package model

import (
	"seedapp/internal/domain/seed"
	"time"
)

type SeedRequest struct {
	TenantType seed.SeedTenantType `json:"tenant_type"`
	Schemas    []string            `json:"schemas"`
	CreatedAt  time.Time           `json:"created_at"`
}

type SeedResponse struct {
	Success []string    `json:"success"`
	Failed  []SeedError `json:"failed"`
}

type SeedError struct {
	Schema string `json:"schema"`
	Error  string `json:"error"`
}
