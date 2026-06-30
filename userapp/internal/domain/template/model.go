package template

import (
	"userapp/internal/domain/shared/entity"
	"userapp/internal/domain/shared/identity"

	"gopkg.in/guregu/null.v4"
)

type Template struct {
	entity.Entity
	Name      string `db:"name" json:"name"`
	Category  string `db:"category" json:"category"`
	Published bool   `db:"published" json:"published"`
}

type CreateTemplateRequest struct {
	Name      string `json:"name"`
	Category  string `json:"category"`
	Published bool   `json:"published"`
}

type UpdateTemplateRequest struct {
	ID        identity.ID `json:"id"`
	Name      string      `json:"name"`
	Category  string      `json:"category"`
	Published null.Bool   `json:"published"`
}
