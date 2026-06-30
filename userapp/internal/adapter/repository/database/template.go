package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"userapp/internal/adapter/repository"
	"userapp/internal/domain/shared/identity"
	"userapp/internal/domain/template"

	"github.com/runsystemid/golog"
)

type TemplateRepository struct {
	DB repository.Sqlx `inject:"database"`
}

func (r *TemplateRepository) Startup() error {
	// Table migrations are now managed globally via golang-migrate in database.go
	return nil
}

func (r *TemplateRepository) Shutdown() error {
	return nil
}

func (t *TemplateRepository) Create(ctx context.Context, data *template.Template) (*template.Template, error) {
	query := `INSERT INTO templates (id, name, category, published, created_at, updated_at, deleted_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := t.DB.GetDB().ExecContext(ctx, query, data.ID.String(), data.Name, data.Category, data.Published, data.CreatedAt, data.UpdatedAt, data.DeletedAt)
	if err != nil {
		if t.DB.IsErrorDuplicate(err) {
			return nil, template.ErrDataAlreadyExists
		}
		golog.Error(ctx, fmt.Sprintf("Error create template: %s", err.Error()), err)
		return nil, template.ErrDataCreate
	}

	return data, nil
}

func (t *TemplateRepository) Update(ctx context.Context, data *template.Template) error {
	query := `UPDATE templates SET name = $1, category = $2, published = $3, updated_at = $4, deleted_at = $5 WHERE id = $6`
	_, err := t.DB.GetDB().ExecContext(ctx, query, data.Name, data.Category, data.Published, data.UpdatedAt, data.DeletedAt, data.ID.String())
	if err != nil {
		golog.Error(ctx, fmt.Sprintf("Error update template: %s", err.Error()), err)
		return template.ErrDataUpdate
	}

	return nil
}

func (t *TemplateRepository) GetByID(ctx context.Context, id identity.ID) (*template.Template, error) {
	query := `SELECT id, name, category, published, created_at, updated_at, deleted_at FROM templates WHERE id = $1 AND deleted_at IS NULL`
	var s template.Template
	err := t.DB.GetDB().GetContext(ctx, &s, query, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, template.ErrDataNotFound
		}
		golog.Error(ctx, fmt.Sprintf("Error get by ID: %s", err.Error()), err)
		return nil, template.ErrDataGet
	}

	return &s, nil
}
