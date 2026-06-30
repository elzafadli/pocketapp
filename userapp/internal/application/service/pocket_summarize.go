package service

import (
	"context"

	"userapp/internal/domain/agent"
	"userapp/internal/domain/pocket"
	"userapp/internal/domain/shared/identity"

	"github.com/runsystemid/golog"
)

func (s *Pocket) Summarize(ctx context.Context, schema string, id identity.ID) (*agent.SummaryResponse, error) {
	item, err := s.PocketRepo.GetByID(ctx, schema, id)
	if err != nil {
		golog.Error(ctx, "failed to find pocket item", err)
		return nil, err
	}
	if item == nil {
		return nil, pocket.ErrPocketNotFound
	}

	res, err := s.AgentClient.Summarize(ctx, item)
	if err != nil {
		golog.Error(ctx, "failed to summarize pocket item via agent", err)
		return nil, err
	}

	return res, nil
}
