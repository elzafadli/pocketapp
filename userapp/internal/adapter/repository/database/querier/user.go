package querier

import (
	"userapp/internal/domain/user"

	sq "github.com/Masterminds/squirrel"
)

const (
	USER_TABLE = "main.users"
)

type UserQuerier interface {
	Create(data *user.User) (string, []interface{}, error)
	Update(data *user.User) (string, []interface{}, error)
	GetByID(id string) (string, []interface{}, error)
	GetByEmail(email string) (string, []interface{}, error)
	Delete(id string) (string, []interface{}, error)
	List(filter map[string]interface{}) (string, []interface{}, error)
	Count(filter map[string]interface{}) (string, []interface{}, error)
}

type User struct {
	SQLBuilder sq.StatementBuilderType `inject:"sqlBuilder"`
}

func (q *User) Create(data *user.User) (string, []interface{}, error) {
	return q.SQLBuilder.
		Insert(USER_TABLE).
		Columns("user_code", "user_name", "tenant_default", "active_indicator", "email", "password", "created_at", "updated_at").
		Values(data.ID, data.Name, data.TenantDefault, data.ActiveIndicator, data.Email, data.Password, data.CreatedAt, data.UpdatedAt).
		ToSql()
}

func (q *User) Update(data *user.User) (string, []interface{}, error) {
	return q.SQLBuilder.
		Update(USER_TABLE).
		Set("user_name", data.Name).
		Set("tenant_default", data.TenantDefault).
		Set("active_indicator", data.ActiveIndicator).
		Set("email", data.Email).
		Set("password", data.Password).
		Set("updated_at", data.UpdatedAt).
		Where(sq.Eq{"user_code": data.ID}).
		ToSql()
}

func (q *User) GetByID(id string) (string, []interface{}, error) {
	return q.SQLBuilder.
		Select("user_code AS id", "user_name AS name", "tenant_default", "active_indicator", "email", "password", "created_at", "updated_at").
		From(USER_TABLE).
		Where(sq.Eq{"user_code": id}).
		ToSql()
}

func (q *User) GetByEmail(email string) (string, []interface{}, error) {
	return q.SQLBuilder.
		Select("user_code AS id", "user_name AS name", "tenant_default", "active_indicator", "email", "password", "created_at", "updated_at").
		From(USER_TABLE).
		Where(sq.Eq{"email": email}).
		ToSql()
}

func (q *User) Delete(id string) (string, []interface{}, error) {
	return q.SQLBuilder.
		Delete(USER_TABLE).
		Where(sq.Eq{"user_code": id}).
		ToSql()
}

func (q *User) List(filter map[string]interface{}) (string, []interface{}, error) {
	sql := q.SQLBuilder.
		Select("user_code AS id", "user_name AS name", "tenant_default", "active_indicator", "email", "password", "created_at", "updated_at").
		From(USER_TABLE)

	if name, ok := filter["name"].(string); ok && name != "" {
		sql = sql.Where(sq.ILike{"user_name": "%" + name + "%"})
	}
	if email, ok := filter["email"].(string); ok && email != "" {
		sql = sql.Where(sq.ILike{"email": "%" + email + "%"})
	}

	sql = sql.OrderBy("created_at DESC")

	if limit, ok := filter["limit"].(int); ok && limit > 0 {
		sql = sql.Limit(uint64(limit))
	}

	if offset, ok := filter["offset"].(int); ok && offset >= 0 {
		sql = sql.Offset(uint64(offset))
	}

	return sql.ToSql()
}

func (q *User) Count(filter map[string]interface{}) (string, []interface{}, error) {
	sql := q.SQLBuilder.
		Select("COUNT(1)").
		From(USER_TABLE)

	if name, ok := filter["name"].(string); ok && name != "" {
		sql = sql.Where(sq.ILike{"user_name": "%" + name + "%"})
	}
	if email, ok := filter["email"].(string); ok && email != "" {
		sql = sql.Where(sq.ILike{"email": "%" + email + "%"})
	}

	return sql.ToSql()
}
