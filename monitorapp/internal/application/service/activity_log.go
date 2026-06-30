package service

import (
	"context"
	"time"

	"monitorapp/internal/domain/activity_log"

	"github.com/runsystemid/golog"
)

type ActivityLogService interface {
	Create(ctx context.Context, data *activity_log.CreateActivityLogRequest) (*activity_log.ActivityLog, error)
	Update(ctx context.Context, data *activity_log.UpdateActivityLogRequest) (*activity_log.ActivityLog, error)
	Find(ctx context.Context, id string) (*activity_log.ActivityLog, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter map[string]interface{}) ([]*activity_log.ActivityLog, uint64, error)
}

type ActivityLog struct {
	ActivityLogRepo activity_log.Repository `inject:"activityLogRepository"`
}

func (s *ActivityLog) Create(ctx context.Context, param *activity_log.CreateActivityLogRequest) (*activity_log.ActivityLog, error) {
	now := time.Now()
	a := activity_log.ActivityLog{
		ObjectName:  param.ObjectName,
		RecordID:    param.RecordID,
		Action:      param.Action,
		ChangedBy:   param.ChangedBy,
		Request:     param.Request,
		Response:    param.Response,
		ChangeStamp: now,
	}

	newLog, err := s.ActivityLogRepo.Create(ctx, &a)
	if err != nil {
		golog.Error(ctx, "Error create activity log: "+err.Error(), err)
		return nil, err
	}

	return newLog, nil
}

func (s *ActivityLog) Update(ctx context.Context, param *activity_log.UpdateActivityLogRequest) (*activity_log.ActivityLog, error) {
	a, err := s.ActivityLogRepo.GetByID(ctx, param.ID)
	if err != nil {
		return nil, err
	}

	if param.ObjectName != "" {
		a.ObjectName = param.ObjectName
	}
	if param.RecordID != "" {
		a.RecordID = param.RecordID
	}
	if param.Action != "" {
		a.Action = param.Action
	}
	if param.ChangedBy != "" {
		a.ChangedBy = param.ChangedBy
	}
	if param.Request != nil {
		a.Request = param.Request
	}
	if param.Response != nil {
		a.Response = param.Response
	}
	a.ChangeStamp = time.Now()

	err = s.ActivityLogRepo.Update(ctx, a)
	if err != nil {
		golog.Error(ctx, "Error update activity log: "+err.Error(), err)
		return nil, err
	}

	return a, nil
}

func (s *ActivityLog) Find(ctx context.Context, id string) (*activity_log.ActivityLog, error) {
	return s.ActivityLogRepo.GetByID(ctx, id)
}

func (s *ActivityLog) Delete(ctx context.Context, id string) error {
	err := s.ActivityLogRepo.Delete(ctx, id)
	if err != nil {
		golog.Error(ctx, "Error delete activity log: "+err.Error(), err)
		return err
	}
	return nil
}

func (s *ActivityLog) List(ctx context.Context, filter map[string]interface{}) ([]*activity_log.ActivityLog, uint64, error) {
	logs, err := s.ActivityLogRepo.GetAll(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.ActivityLogRepo.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
