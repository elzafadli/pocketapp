package service

import (
	"context"
	"time"

	"userapp/internal/domain/agent"
	"userapp/internal/domain/pocket"
	"userapp/internal/domain/shared/identity"

	"github.com/runsystemid/golog"
)

type PocketService interface {
	Create(ctx context.Context, schema string, data *pocket.CreatePocketRequest) (*pocket.PocketItem, error)
	Update(ctx context.Context, schema string, data *pocket.UpdatePocketRequest) (*pocket.PocketItem, error)
	Find(ctx context.Context, schema string, id identity.ID) (*pocket.PocketItem, error)
	Delete(ctx context.Context, schema string, id identity.ID) error
	List(ctx context.Context, schema string, query *pocket.PocketListQuery) ([]*pocket.PocketItem, uint64, error)
	UpdateStatus(ctx context.Context, schema string, id identity.ID, status string) (*pocket.PocketItem, error)
	ToggleFavorite(ctx context.Context, schema string, id identity.ID, favorite bool) (*pocket.PocketItem, error)
	Summarize(ctx context.Context, schema string, id identity.ID) (*agent.SummaryResponse, error)
}

type Pocket struct {
	PocketRepo  pocket.Repository `inject:"pocketRepository"`
	AgentClient agent.AgentClient `inject:"agentClient"`
}

func (s *Pocket) Create(ctx context.Context, schema string, param *pocket.CreatePocketRequest) (*pocket.PocketItem, error) {
	now := time.Now()

	var urlPtr *string
	if param.URL != "" {
		urlVal := param.URL
		urlPtr = &urlVal
	}

	var descPtr *string
	if param.Description != "" {
		descVal := param.Description
		descPtr = &descVal
	}

	tags := pocket.Tags(param.Tags)
	if tags == nil {
		tags = []string{}
	}

	item := &pocket.PocketItem{
		ID:          identity.NewID(),
		Title:       param.Title,
		URL:         urlPtr,
		Description: descPtr,
		ContentType: param.ContentType,
		Status:      "unread",
		IsFavorite:  false,
		Tags:        tags,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	newItem, err := s.PocketRepo.Create(ctx, schema, item)
	if err != nil {
		golog.Error(ctx, "Error create pocket: "+err.Error(), err)
		return nil, err
	}

	return newItem, nil
}

func (s *Pocket) Update(ctx context.Context, schema string, param *pocket.UpdatePocketRequest) (*pocket.PocketItem, error) {
	item, err := s.PocketRepo.GetByID(ctx, schema, param.ID)
	if err != nil {
		return nil, err
	}

	item.Title = param.Title

	if param.URL != "" {
		urlVal := param.URL
		item.URL = &urlVal
	} else {
		item.URL = nil
	}

	if param.Description != "" {
		descVal := param.Description
		item.Description = &descVal
	} else {
		item.Description = nil
	}

	item.ContentType = param.ContentType
	item.Tags = pocket.Tags(param.Tags)
	if item.Tags == nil {
		item.Tags = []string{}
	}
	item.UpdatedAt = time.Now()

	err = s.PocketRepo.Update(ctx, schema, item)
	if err != nil {
		golog.Error(ctx, "Error update pocket: "+err.Error(), err)
		return nil, err
	}

	return item, nil
}

func (s *Pocket) Find(ctx context.Context, schema string, id identity.ID) (*pocket.PocketItem, error) {
	return s.PocketRepo.GetByID(ctx, schema, id)
}

func (s *Pocket) Delete(ctx context.Context, schema string, id identity.ID) error {
	return s.PocketRepo.Delete(ctx, schema, id)
}

func (s *Pocket) List(ctx context.Context, schema string, query *pocket.PocketListQuery) ([]*pocket.PocketItem, uint64, error) {
	filter := query.ToFilter()

	list, err := s.PocketRepo.List(ctx, schema, filter)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.PocketRepo.Count(ctx, schema, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (s *Pocket) UpdateStatus(ctx context.Context, schema string, id identity.ID, status string) (*pocket.PocketItem, error) {
	item, err := s.PocketRepo.GetByID(ctx, schema, id)
	if err != nil {
		return nil, err
	}

	item.Status = status
	item.UpdatedAt = time.Now()

	err = s.PocketRepo.Update(ctx, schema, item)
	if err != nil {
		golog.Error(ctx, "Error update status pocket: "+err.Error(), err)
		return nil, err
	}

	return item, nil
}

func (s *Pocket) ToggleFavorite(ctx context.Context, schema string, id identity.ID, favorite bool) (*pocket.PocketItem, error) {
	item, err := s.PocketRepo.GetByID(ctx, schema, id)
	if err != nil {
		return nil, err
	}

	item.IsFavorite = favorite
	item.UpdatedAt = time.Now()

	err = s.PocketRepo.Update(ctx, schema, item)
	if err != nil {
		golog.Error(ctx, "Error toggle favorite pocket: "+err.Error(), err)
		return nil, err
	}

	return item, nil
}
