package activity_log

import (
	"errors"
	"time"
)

type ObjectName string
type Action string

type ActivityLog struct {
	ID          string      `db:"id" json:"id"`
	ObjectName  ObjectName  `db:"object_name" json:"object_name"`
	RecordID    string      `db:"record_id" json:"record_id"`
	Action      Action      `db:"action" json:"action"`
	ChangedBy   string      `db:"changed_by" json:"changed_by"`
	Request     interface{} `db:"request" json:"request"`
	Response    interface{} `db:"response" json:"response"`
	ChangeStamp time.Time   `db:"change_stamp" json:"change_stamp"`
}

type CreateActivityLogRequest struct {
	ObjectName  ObjectName  `json:"object_name" validate:"required,max=100"`
	RecordID    string      `json:"record_id" validate:"required,max=100"`
	Action      Action      `json:"action" validate:"required,max=100"`
	ChangedBy   string      `json:"changed_by" validate:"required,max=255"`
	Request     interface{} `json:"request"`
	Response    interface{} `json:"response"`
}

type UpdateActivityLogRequest struct {
	ID          string      `json:"id" validate:"required"`
	ObjectName  ObjectName  `json:"object_name" validate:"max=100"`
	RecordID    string      `json:"record_id" validate:"max=100"`
	Action      Action      `json:"action" validate:"max=100"`
	ChangedBy   string      `json:"changed_by" validate:"max=255"`
	Request     interface{} `json:"request"`
	Response    interface{} `json:"response"`
}

var (
	ErrDataNotFound      = errors.New("activity log data not found")
	ErrDataAlreadyExists = errors.New("activity log data already exists")
	ErrDataCreate        = errors.New("failed to create activity log")
	ErrDataUpdate        = errors.New("failed to update activity log")
	ErrDataDelete        = errors.New("failed to delete activity log")
	ErrDataGet           = errors.New("failed to get activity log")
)
