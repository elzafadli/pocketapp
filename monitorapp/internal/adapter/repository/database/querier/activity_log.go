package querier

import (
	"monitorapp/internal/domain/activity_log"

	sq "github.com/Masterminds/squirrel"
)

const (
	ACTIVITY_LOG_TABLE = "user_activity_logs"
)

type ActivityLogQuerier interface {
	Create(data *activity_log.ActivityLog, requestJSON, responseJSON []byte) (string, []interface{}, error)
	Update(data *activity_log.ActivityLog, requestJSON, responseJSON []byte) (string, []interface{}, error)
	GetByID(id string) (string, []interface{}, error)
	Delete(id string) (string, []interface{}, error)
	List(filter map[string]interface{}) (string, []interface{}, error)
	Count(filter map[string]interface{}) (string, []interface{}, error)
}

type ActivityLog struct {
	SQLBuilder sq.StatementBuilderType `inject:"sqlBuilder"`
}

func (q *ActivityLog) Create(data *activity_log.ActivityLog, requestJSON, responseJSON []byte) (string, []interface{}, error) {
	return q.SQLBuilder.
		Insert(ACTIVITY_LOG_TABLE).
		Columns("id", "object_name", "record_id", "action", "changed_by", "request", "response", "change_stamp").
		Values(data.ID, string(data.ObjectName), data.RecordID, string(data.Action), data.ChangedBy, requestJSON, responseJSON, data.ChangeStamp).
		ToSql()
}

func (q *ActivityLog) Update(data *activity_log.ActivityLog, requestJSON, responseJSON []byte) (string, []interface{}, error) {
	return q.SQLBuilder.
		Update(ACTIVITY_LOG_TABLE).
		Set("object_name", string(data.ObjectName)).
		Set("record_id", data.RecordID).
		Set("action", string(data.Action)).
		Set("changed_by", data.ChangedBy).
		Set("request", requestJSON).
		Set("response", responseJSON).
		Set("change_stamp", data.ChangeStamp).
		Where(sq.Eq{"id": data.ID}).
		ToSql()
}

func (q *ActivityLog) GetByID(id string) (string, []interface{}, error) {
	return q.SQLBuilder.
		Select("id", "object_name", "record_id", "action", "changed_by", "request", "response", "change_stamp").
		From(ACTIVITY_LOG_TABLE).
		Where(sq.Eq{"id": id}).
		ToSql()
}

func (q *ActivityLog) Delete(id string) (string, []interface{}, error) {
	return q.SQLBuilder.
		Delete(ACTIVITY_LOG_TABLE).
		Where(sq.Eq{"id": id}).
		ToSql()
}

func (q *ActivityLog) List(filter map[string]interface{}) (string, []interface{}, error) {
	sql := q.SQLBuilder.
		Select("id", "object_name", "record_id", "action", "changed_by", "request", "response", "change_stamp").
		From(ACTIVITY_LOG_TABLE)

	if objectName, ok := filter["object_name"].(string); ok && objectName != "" {
		sql = sql.Where(sq.Eq{"object_name": objectName})
	}
	if recordID, ok := filter["record_id"].(string); ok && recordID != "" {
		sql = sql.Where(sq.Eq{"record_id": recordID})
	}
	if action, ok := filter["action"].(string); ok && action != "" {
		sql = sql.Where(sq.Eq{"action": action})
	}
	if changedBy, ok := filter["changed_by"].(string); ok && changedBy != "" {
		sql = sql.Where(sq.ILike{"changed_by": "%" + changedBy + "%"})
	}

	sql = sql.OrderBy("change_stamp DESC")

	if limit, ok := filter["limit"].(int); ok && limit > 0 {
		sql = sql.Limit(uint64(limit))
	}

	if offset, ok := filter["offset"].(int); ok && offset >= 0 {
		sql = sql.Offset(uint64(offset))
	}

	return sql.ToSql()
}

func (q *ActivityLog) Count(filter map[string]interface{}) (string, []interface{}, error) {
	sql := q.SQLBuilder.
		Select("COUNT(1)").
		From(ACTIVITY_LOG_TABLE)

	if objectName, ok := filter["object_name"].(string); ok && objectName != "" {
		sql = sql.Where(sq.Eq{"object_name": objectName})
	}
	if recordID, ok := filter["record_id"].(string); ok && recordID != "" {
		sql = sql.Where(sq.Eq{"record_id": recordID})
	}
	if action, ok := filter["action"].(string); ok && action != "" {
		sql = sql.Where(sq.Eq{"action": action})
	}
	if changedBy, ok := filter["changed_by"].(string); ok && changedBy != "" {
		sql = sql.Where(sq.ILike{"changed_by": "%" + changedBy + "%"})
	}

	return sql.ToSql()
}
