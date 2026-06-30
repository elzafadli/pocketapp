package agent

import (
	"context"

	"userapp/internal/domain/pocket"
)

type SummaryResponse struct {
	Summary  string   `json:"summary"`
	TodoList []string `json:"todoList"`
}

type AgentClient interface {
	Summarize(ctx context.Context, item *pocket.PocketItem) (*SummaryResponse, error)
}
