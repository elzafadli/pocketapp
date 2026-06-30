package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"userapp/internal/adapter/repository"
	"userapp/internal/adapter/repository/database/querier"
	"userapp/internal/domain/pocket"
	"userapp/internal/domain/shared/identity"

	"github.com/runsystemid/golog"
)

type PocketRepository struct {
	DB      repository.Sqlx       `inject:"database"`
	Querier querier.PocketQuerier `inject:"pocketQuerier"`
}

func (r *PocketRepository) Startup() error  { return nil }
func (r *PocketRepository) Shutdown() error { return nil }

func (r *PocketRepository) Create(ctx context.Context, schema string, data *pocket.PocketItem) (*pocket.PocketItem, error) {
	query, args, err := r.Querier.Create(schema, data)
	if err != nil {
		golog.Error(ctx, fmt.Sprintf("Error build query create pocket: %s", err.Error()), err)
		return nil, pocket.ErrDataCreate
	}

	_, err = r.DB.ExecContext(ctx, query, args...)
	if err != nil {
		if r.DB.IsErrorDuplicate(err) {
			return nil, pocket.ErrDataAlreadyExists
		}
		golog.Error(ctx, fmt.Sprintf("Error create pocket: %s", err.Error()), err)
		return nil, pocket.ErrDataCreate
	}

	return data, nil
}

func (r *PocketRepository) Update(ctx context.Context, schema string, data *pocket.PocketItem) error {
	query, args, err := r.Querier.Update(schema, data)
	if err != nil {
		golog.Error(ctx, fmt.Sprintf("Error build query update pocket: %s", err.Error()), err)
		return pocket.ErrDataUpdate
	}

	result, err := r.DB.GetDB().ExecContext(ctx, query, args...)
	if err != nil {
		if r.DB.IsErrorDuplicate(err) {
			return pocket.ErrDataAlreadyExists
		}
		golog.Error(ctx, fmt.Sprintf("Error update pocket: %s", err.Error()), err)
		return pocket.ErrDataUpdate
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		golog.Error(ctx, fmt.Sprintf("Error getting rows affected in update pocket: %s", err.Error()), err)
		return pocket.ErrDataUpdate
	}

	if rowsAffected == 0 {
		return pocket.ErrPocketNotFound
	}

	return nil
}

func (r *PocketRepository) GetByID(ctx context.Context, schema string, id identity.ID) (*pocket.PocketItem, error) {
	query, args, err := r.Querier.GetByID(schema, id.String())
	if err != nil {
		golog.Error(ctx, fmt.Sprintf("Error build query get pocket by ID: %s", err.Error()), err)
		return nil, pocket.ErrDataGet
	}

	var data pocket.PocketItem
	err = r.DB.GetDB().GetContext(ctx, &data, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pocket.ErrPocketNotFound
		}
		golog.Error(ctx, fmt.Sprintf("Error get pocket by ID: %s", err.Error()), err)
		return nil, pocket.ErrDataGet
	}

	return &data, nil
}

func (r *PocketRepository) Delete(ctx context.Context, schema string, id identity.ID) error {
	now := time.Now()
	query, args, err := r.Querier.Delete(schema, id.String(), now)
	if err != nil {
		golog.Error(ctx, fmt.Sprintf("Error build query delete pocket: %s", err.Error()), err)
		return pocket.ErrDataDelete
	}

	result, err := r.DB.GetDB().ExecContext(ctx, query, args...)
	if err != nil {
		golog.Error(ctx, fmt.Sprintf("Error delete pocket: %s", err.Error()), err)
		return pocket.ErrDataDelete
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		golog.Error(ctx, fmt.Sprintf("Error getting rows affected in delete pocket: %s", err.Error()), err)
		return pocket.ErrDataDelete
	}

	if rowsAffected == 0 {
		return pocket.ErrPocketNotFound
	}

	return nil
}

func (r *PocketRepository) List(ctx context.Context, schema string, filter map[string]interface{}) ([]*pocket.PocketItem, error) {
	query, args, err := r.Querier.List(schema, filter)
	if err != nil {
		golog.Error(ctx, fmt.Sprintf("Error build query list pocket: %s", err.Error()), err)
		return nil, pocket.ErrDataGet
	}

	var list []*pocket.PocketItem
	err = r.DB.GetDB().SelectContext(ctx, &list, query, args...)
	if err != nil {
		golog.Error(ctx, fmt.Sprintf("Error select list pocket: %s", err.Error()), err)
		return nil, pocket.ErrDataGet
	}

	return list, nil
}

func (r *PocketRepository) Count(ctx context.Context, schema string, filter map[string]interface{}) (uint64, error) {
	query, args, err := r.Querier.Count(schema, filter)
	if err != nil {
		golog.Error(ctx, fmt.Sprintf("Error build query count pocket: %s", err.Error()), err)
		return 0, pocket.ErrDataGet
	}

	var count uint64
	err = r.DB.GetDB().GetContext(ctx, &count, query, args...)
	if err != nil {
		golog.Error(ctx, fmt.Sprintf("Error get count pocket: %s", err.Error()), err)
		return 0, pocket.ErrDataGet
	}

	return count, nil
}
