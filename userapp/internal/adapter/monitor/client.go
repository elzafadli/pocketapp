package monitor

import (
	"context"
	"errors"
	"fmt"
	"userapp/config"

	"github.com/go-resty/resty/v2"
	"github.com/runsystemid/golog"
)

var ErrFailedToSendActivityLog = errors.New("failed to send activity log")

type MonitorService interface {
	SendActivityLog(ctx context.Context, request *ActivityLogRequest) error
}

type Monitor struct {
	Conf  *config.Config `inject:"config"`
	Resty *resty.Client
}

type ActivityLogRequest struct {
	ObjectName string      `json:"object_name"`
	RecordID   string      `json:"record_id"`
	Action     string      `json:"action"`
	ChangedBy  string      `json:"changed_by"`
	Request    interface{} `json:"request"`
	Response   interface{} `json:"response"`
}

type ActivityLogResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (m *Monitor) Startup() error {
	client := resty.New()
	client.SetBaseURL(m.Conf.Monitor.URL).SetTimeout(m.Conf.Monitor.Timeout)
	client.SetBasicAuth(m.Conf.Monitor.User, m.Conf.Monitor.Password)

	m.Resty = client
	return nil
}

func (m *Monitor) Shutdown() error {
	return nil
}

func (m *Monitor) SendActivityLog(ctx context.Context, request *ActivityLogRequest) error {
	resp, err := m.Resty.R().
		SetContext(ctx).
		SetBody(request).
		SetResult(&ActivityLogResponse{}).
		Post("/v1/activity-log")

	if err != nil {
		golog.Error(ctx, "failed to send activity log: "+err.Error(), err)
		return fmt.Errorf("%w: %s", ErrFailedToSendActivityLog, err.Error())
	}

	if resp.IsError() {
		golog.Error(ctx, fmt.Sprintf("failed to send activity log status code: %d, response: %s", resp.StatusCode(), resp.String()), nil)
		return fmt.Errorf("%w: status %d", ErrFailedToSendActivityLog, resp.StatusCode())
	}

	return nil
}
