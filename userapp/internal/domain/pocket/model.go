package pocket

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"userapp/internal/domain/shared"
	"userapp/internal/domain/shared/identity"
)

type Tags []string

func (t Tags) Value() (driver.Value, error) {
	if t == nil {
		return "[]", nil
	}
	return json.Marshal(t)
}

func (t *Tags) Scan(value interface{}) error {
	if value == nil {
		*t = []string{}
		return nil
	}
	var b []byte
	switch v := value.(type) {
	case string:
		b = []byte(v)
	case []byte:
		b = v
	default:
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, t)
}

type PocketItem struct {
	ID          identity.ID `db:"id" json:"id"`
	Title       string      `db:"title" json:"title"`
	URL         *string     `db:"url" json:"url"`
	Description *string     `db:"description" json:"description"`
	ContentType string      `db:"content_type" json:"contentType"`
	Status      string      `db:"status" json:"status"`
	IsFavorite  bool        `db:"is_favorite" json:"isFavorite"`
	Tags        Tags        `db:"tags" json:"tags"`
	CreatedAt   time.Time   `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time   `db:"updated_at" json:"updatedAt"`
	ArchivedAt  *time.Time  `db:"archived_at" json:"-"`
}

type CreatePocketRequest struct {
	Title       string   `json:"title"        validate:"required,min=3,max=120"`
	URL         string   `json:"url"`
	Description string   `json:"description"  validate:"omitempty,max=500"`
	ContentType string   `json:"contentType"  validate:"required,oneof=article video document note"`
	Tags        []string `json:"tags"         validate:"omitempty,max=10,dive,max=24"`
}

type UpdatePocketRequest struct {
	ID          identity.ID `json:"id"`
	Title       string      `json:"title"        validate:"required,min=3,max=120"`
	URL         string      `json:"url"`
	Description string      `json:"description"  validate:"omitempty,max=500"`
	ContentType string      `json:"contentType"  validate:"required,oneof=article video document note"`
	Tags        []string    `json:"tags"         validate:"omitempty,max=10,dive,max=24"`
}

type UpdateStatusRequest struct {
	Status string `json:"status"`
}

type ToggleFavoriteRequest struct {
	IsFavorite bool `json:"isFavorite"`
}


type PocketListResponse struct {
	Data []*PocketItem       `json:"data"`
	Meta shared.MetaResponse `json:"meta"`
}

type PocketListQuery struct {
	Search      string
	Status      string
	ContentType string
	Favorite    *bool
	Page        int
	Limit       int
	Sort        string
}

func (q *PocketListQuery) ToFilter() map[string]interface{} {
	filter := map[string]interface{}{
		"search": q.Search,
		"status": q.Status,
		"type":   q.ContentType,
		"page":   q.Page,
		"limit":  q.Limit,
		"offset": (q.Page - 1) * q.Limit,
		"sort":   q.Sort,
	}
	if q.Favorite != nil {
		filter["favorite"] = *q.Favorite
	}
	return filter
}
