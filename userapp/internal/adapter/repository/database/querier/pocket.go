package querier

import (
	"userapp/internal/domain/pocket"

	sq "github.com/Masterminds/squirrel"
)

const (
	POCKET_ITEM_TABLE = "pocket_items"
)

type PocketQuerier interface {
	Create(schema string, data *pocket.PocketItem) (string, []interface{}, error)
	Update(schema string, data *pocket.PocketItem) (string, []interface{}, error)
	GetByID(schema string, id string) (string, []interface{}, error)
	Delete(schema string, id string, archivedAt interface{}) (string, []interface{}, error)
	List(schema string, filter map[string]interface{}) (string, []interface{}, error)
	Count(schema string, filter map[string]interface{}) (string, []interface{}, error)
}

type Pocket struct {
	SQLBuilder sq.StatementBuilderType `inject:"sqlBuilder"`
}

func (q *Pocket) Create(schema string, data *pocket.PocketItem) (string, []interface{}, error) {
	return q.SQLBuilder.
		Insert(TenantTableSchema(schema, POCKET_ITEM_TABLE)).
		Columns("id", "title", "url", "description", "content_type", "status", "is_favorite", "tags", "created_at", "updated_at", "archived_at").
		Values(data.ID.String(), data.Title, data.URL, data.Description, data.ContentType, data.Status, data.IsFavorite, data.Tags, data.CreatedAt, data.UpdatedAt, data.ArchivedAt).
		ToSql()
}

func (q *Pocket) Update(schema string, data *pocket.PocketItem) (string, []interface{}, error) {
	return q.SQLBuilder.
		Update(TenantTableSchema(schema, POCKET_ITEM_TABLE)).
		Set("title", data.Title).
		Set("url", data.URL).
		Set("description", data.Description).
		Set("content_type", data.ContentType).
		Set("status", data.Status).
		Set("is_favorite", data.IsFavorite).
		Set("tags", data.Tags).
		Set("updated_at", data.UpdatedAt).
		Set("archived_at", data.ArchivedAt).
		Where(sq.Eq{"id": data.ID.String()}).
		ToSql()
}

func (q *Pocket) GetByID(schema string, id string) (string, []interface{}, error) {
	return q.SQLBuilder.
		Select("id", "title", "url", "description", "content_type", "status", "is_favorite", "tags", "created_at", "updated_at", "archived_at").
		From(TenantTableSchema(schema, POCKET_ITEM_TABLE)).
		Where(sq.Eq{"id": id}).
		ToSql()
}

func (q *Pocket) Delete(schema string, id string, archivedAt interface{}) (string, []interface{}, error) {
	return q.SQLBuilder.
		Update(TenantTableSchema(schema, POCKET_ITEM_TABLE)).
		Set("archived_at", archivedAt).
		Set("status", "archived").
		Set("updated_at", archivedAt).
		Where(sq.Eq{"id": id}).
		ToSql()
}

func (q *Pocket) List(schema string, filter map[string]interface{}) (string, []interface{}, error) {
	sql := q.SQLBuilder.
		Select("id", "title", "url", "description", "content_type", "status", "is_favorite", "tags", "created_at", "updated_at", "archived_at").
		From(TenantTableSchema(schema, POCKET_ITEM_TABLE))

	sql = q.applyFilters(sql, filter)

	// Sort logic
	sort := "created_at DESC"
	if sortVal, ok := filter["sort"].(string); ok && sortVal != "" {
		if sortVal == "createdAt:asc" {
			sort = "created_at ASC"
		} else if sortVal == "createdAt:desc" {
			sort = "created_at DESC"
		}
	}
	sql = sql.OrderBy(sort)

	if limit, ok := filter["limit"].(int); ok && limit > 0 {
		sql = sql.Limit(uint64(limit))
	}

	if offset, ok := filter["offset"].(int); ok && offset >= 0 {
		sql = sql.Offset(uint64(offset))
	}

	return sql.ToSql()
}

func (q *Pocket) Count(schema string, filter map[string]interface{}) (string, []interface{}, error) {
	sql := q.SQLBuilder.
		Select("COUNT(1)").
		From(TenantTableSchema(schema, POCKET_ITEM_TABLE))

	sql = q.applyFilters(sql, filter)

	return sql.ToSql()
}

func (q *Pocket) applyFilters(sql sq.SelectBuilder, filter map[string]interface{}) sq.SelectBuilder {
	// Status filter: defaults to non-archived items if empty or not "archived"
	status, hasStatus := filter["status"].(string)
	if hasStatus && status != "" {
		if status == "archived" {
			sql = sql.Where(sq.NotEq{"archived_at": nil})
			sql = sql.Where(sq.Eq{"status": "archived"})
		} else {
			sql = sql.Where(sq.Eq{"archived_at": nil})
			sql = sql.Where(sq.Eq{"status": status})
		}
	} else {
		// By default list only active (non-archived) items
		sql = sql.Where(sq.Eq{"archived_at": nil})
	}

	// Search filter
	if search, ok := filter["search"].(string); ok && search != "" {
		sql = sql.Where(sq.Or{
			sq.ILike{"title": "%" + search + "%"},
			sq.ILike{"description": "%" + search + "%"},
		})
	}

	// Content Type filter
	if contentType, ok := filter["type"].(string); ok && contentType != "" {
		sql = sql.Where(sq.Eq{"content_type": contentType})
	}

	// Favorite filter
	if favorite, ok := filter["favorite"].(bool); ok {
		sql = sql.Where(sq.Eq{"is_favorite": favorite})
	}

	return sql
}
