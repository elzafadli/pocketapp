package service

import (
	"context"
	"fmt"

	"github.com/runsystemid/golog"
	"userapp/internal/adapter/monitor"
)

type ActivityLogRequest struct {
	ObjectName string      `json:"object_name"`
	RecordID   string      `json:"record_id"`
	Action     string      `json:"action"`
	ChangedBy  string      `json:"changed_by"`
	Request    interface{} `json:"request"`
	Response   interface{} `json:"response"`
}

//go:generate mockgen -destination=mocks/activity_log.go -package=mocks -source=activity_log.go
type ActivityLogService interface {
	Send(ctx context.Context, req *ActivityLogRequest) error
}

type ActivityLog struct {
	MonitorClient monitor.MonitorService `inject:"monitorClient"`
}

func (s *ActivityLog) Send(ctx context.Context, req *ActivityLogRequest) error {
	payload := &monitor.ActivityLogRequest{
		ObjectName: req.ObjectName,
		RecordID:   req.RecordID,
		Action:     req.Action,
		ChangedBy:  req.ChangedBy,
		Request:    req.Request,
		Response:   req.Response,
	}

	err := s.MonitorClient.SendActivityLog(ctx, payload)
	if err != nil {
		golog.Error(ctx, "failed to send activity log: "+err.Error(), err)
		return fmt.Errorf("failed to send activity log: %w", err)
	}

	return nil
}
