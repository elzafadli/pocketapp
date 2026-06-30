package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"monitorapp/internal/adapter/repository"
	"monitorapp/internal/adapter/repository/database/querier"
	"monitorapp/internal/domain/activity_log"

	"github.com/google/uuid"
	"github.com/runsystemid/golog"
)

type activityLogDB struct {
	ID          string    `db:"id"`
	ObjectName  string    `db:"object_name"`
	RecordID    string    `db:"record_id"`
	Action      string    `db:"action"`
	ChangedBy   string    `db:"changed_by"`
	Request     []byte    `db:"request"`
	Response    []byte    `db:"response"`
	ChangeStamp time.Time `db:"change_stamp"`
}

type ActivityLogRepository struct {
	DB      repository.Sqlx            `inject:"database"`
	Querier querier.ActivityLogQuerier `inject:"activityLogQuerier"`
}

func (r *ActivityLogRepository) Startup() error {
	return nil
}

func (r *ActivityLogRepository) Shutdown() error {
	return nil
}

func (r *ActivityLogRepository) Create(ctx context.Context, data *activity_log.ActivityLog) (*activity_log.ActivityLog, error) {
	if data.ID == "" {
		data.ID = uuid.New().String()
	}

	var reqJSON, resJSON []byte
	var err error
	if data.Request != nil {
		reqJSON, err = json.Marshal(data.Request)
		if err != nil {
			golog.Error(ctx, fmt.Sprintf("Error marshal request: %s", err.Error()), err)
			return nil, activity_log.ErrDataCreate
		}
	}
	if data.Response != nil {
		resJSON, err = json.Marshal(data.Response)
		if err != nil {
			golog.Error(ctx, fmt.Sprintf("Error marshal response: %s", err.Error()), err)
			return nil, activity_log.ErrDataCreate
		}
	}

	query, args, err := r.Querier.Create(data, reqJSON, resJSON)
	if err != nil {
		golog.Error(ctx, fmt.Sprintf("Error build query create activity log: %s", err.Error()), err)
		return nil, activity_log.ErrDataCreate
	}

	_, err = r.DB.ExecContext(ctx, query, args...)
	if err != nil {
		if r.DB.IsErrorDuplicate(err) {
			return nil, activity_log.ErrDataAlreadyExists
		}
		golog.Error(ctx, fmt.Sprintf("Error create activity log: %s", err.Error()), err)
		return nil, activity_log.ErrDataCreate
	}

	return data, nil
}

func (r *ActivityLogRepository) Update(ctx context.Context, data *activity_log.ActivityLog) error {
	var reqJSON, resJSON []byte
	var err error
	if data.Request != nil {
		reqJSON, err = json.Marshal(data.Request)
		if err != nil {
			golog.Error(ctx, fmt.Sprintf("Error marshal request in update: %s", err.Error()), err)
			return activity_log.ErrDataUpdate
		}
	}
	if data.Response != nil {
		resJSON, err = json.Marshal(data.Response)
		if err != nil {
			golog.Error(ctx, fmt.Sprintf("Error marshal response in update: %s", err.Error()), err)
			return activity_log.ErrDataUpdate
		}
	}

	query, args, err := r.Querier.Update(data, reqJSON, resJSON)
	if err != nil {
		golog.Error(ctx, fmt.Sprintf("Error build query update activity log: %s", err.Error()), err)
		return activity_log.ErrDataUpdate
	}

	result, err := r.DB.GetDB().ExecContext(ctx, query, args...)
	if err != nil {
		if r.DB.IsErrorDuplicate(err) {
			return activity_log.ErrDataAlreadyExists
		}
		golog.Error(ctx, fmt.Sprintf("Error update activity log: %s", err.Error()), err)
		return activity_log.ErrDataUpdate
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		golog.Error(ctx, fmt.Sprintf("Error getting rows affected in update activity log: %s", err.Error()), err)
		return activity_log.ErrDataUpdate
	}

	if rowsAffected == 0 {
		return activity_log.ErrDataNotFound
	}

	return nil
}

func (r *ActivityLogRepository) GetByID(ctx context.Context, id string) (*activity_log.ActivityLog, error) {
	query, args, err := r.Querier.GetByID(id)
	if err != nil {
		golog.Error(ctx, fmt.Sprintf("Error build query get activity log by ID: %s", err.Error()), err)
		return nil, activity_log.ErrDataGet
	}

	var dbLog activityLogDB
	err = r.DB.GetDB().GetContext(ctx, &dbLog, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, activity_log.ErrDataNotFound
		}
		golog.Error(ctx, fmt.Sprintf("Error get activity log by ID: %s", err.Error()), err)
		return nil, activity_log.ErrDataGet
	}

	return mapToDomain(&dbLog), nil
}

func (r *ActivityLogRepository) Delete(ctx context.Context, id string) error {
	query, args, err := r.Querier.Delete(id)
	if err != nil {
		golog.Error(ctx, fmt.Sprintf("Error build query delete activity log: %s", err.Error()), err)
		return activity_log.ErrDataDelete
	}

	result, err := r.DB.GetDB().ExecContext(ctx, query, args...)
	if err != nil {
		golog.Error(ctx, fmt.Sprintf("Error delete activity log: %s", err.Error()), err)
		return activity_log.ErrDataDelete
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		golog.Error(ctx, fmt.Sprintf("Error getting rows affected in delete activity log: %s", err.Error()), err)
		return activity_log.ErrDataDelete
	}

	if rowsAffected == 0 {
		return activity_log.ErrDataNotFound
	}

	return nil
}

func (r *ActivityLogRepository) GetAll(ctx context.Context, filter map[string]interface{}) ([]*activity_log.ActivityLog, error) {
	query, args, err := r.Querier.List(filter)
	if err != nil {
		golog.Error(ctx, fmt.Sprintf("Error build query get all activity logs: %s", err.Error()), err)
		return nil, activity_log.ErrDataGet
	}

	dbLogs := make([]*activityLogDB, 0)
	err = r.DB.GetDB().SelectContext(ctx, &dbLogs, query, args...)
	if err != nil {
		golog.Error(ctx, fmt.Sprintf("Error get all activity logs: %s", err.Error()), err)
		return nil, activity_log.ErrDataGet
	}

	res := make([]*activity_log.ActivityLog, len(dbLogs))
	for i, dbLog := range dbLogs {
		res[i] = mapToDomain(dbLog)
	}

	return res, nil
}

func (r *ActivityLogRepository) Count(ctx context.Context, filter map[string]interface{}) (uint64, error) {
	query, args, err := r.Querier.Count(filter)
	if err != nil {
		golog.Error(ctx, fmt.Sprintf("Error build query count activity logs: %s", err.Error()), err)
		return 0, activity_log.ErrDataGet
	}

	var count uint64
	err = r.DB.GetDB().GetContext(ctx, &count, query, args...)
	if err != nil {
		golog.Error(ctx, fmt.Sprintf("Error count activity logs: %s", err.Error()), err)
		return 0, activity_log.ErrDataGet
	}

	return count, nil
}

func mapToDomain(dbLog *activityLogDB) *activity_log.ActivityLog {
	log := &activity_log.ActivityLog{
		ID:          dbLog.ID,
		ObjectName:  activity_log.ObjectName(dbLog.ObjectName),
		RecordID:    dbLog.RecordID,
		Action:      activity_log.Action(dbLog.Action),
		ChangedBy:   dbLog.ChangedBy,
		ChangeStamp: dbLog.ChangeStamp,
	}

	if len(dbLog.Request) > 0 {
		var req interface{}
		_ = json.Unmarshal(dbLog.Request, &req)
		log.Request = req
	}

	if len(dbLog.Response) > 0 {
		var res interface{}
		_ = json.Unmarshal(dbLog.Response, &res)
		log.Response = res
	}

	return log
}
