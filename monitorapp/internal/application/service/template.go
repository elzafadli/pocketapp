package service

import (
	"context"
	"time"

	"monitorapp/internal/domain/shared/entity"
	"monitorapp/internal/domain/shared/identity"
	"monitorapp/internal/domain/template"

	"github.com/runsystemid/golog"
)

//go:generate mockgen -destination=mocks/template.go -package=mocks -source=template.go
type TemplateService interface {
	Create(ctx context.Context, data *template.CreateTemplateRequest) (*template.Template, error)
	Update(ctx context.Context, data *template.UpdateTemplateRequest) (*template.Template, error)
	Find(ctx context.Context, id identity.ID) (*template.Template, error)
}

type Template struct {
	TemplateRepo template.Repository `inject:"templateCacheRepository"`
}

func (t *Template) Create(ctx context.Context, param *template.CreateTemplateRequest) (*template.Template, error) {
	tem := template.Template{
		Entity:    entity.NewEntity(),
		Name:      param.Name,
		Category:  param.Category,
		Published: param.Published,
	}

	newTemplate, err := t.TemplateRepo.Create(ctx, &tem)
	if err != nil {
		golog.Error(ctx, "Error create template: "+err.Error(), err)
		return nil, err
	}

	return newTemplate, nil
}

func (t *Template) Update(ctx context.Context, param *template.UpdateTemplateRequest) (*template.Template, error) {
	tem, err := t.TemplateRepo.GetByID(ctx, param.ID)
	if err != nil {
		return nil, err
	}

	if param.Name != "" {
		tem.Name = param.Name
	}
	if param.Category != "" {
		tem.Category = param.Category
	}
	if !param.Published.IsZero() {
		tem.Published = param.Published.Bool
	}
	tem.UpdatedAt = time.Now()

	err = t.TemplateRepo.Update(ctx, tem)
	if err != nil {
		golog.Error(ctx, "Error update template: "+err.Error(), err)
		return nil, err
	}

	return tem, nil
}

func (t *Template) Find(ctx context.Context, id identity.ID) (*template.Template, error) {
	return t.TemplateRepo.GetByID(ctx, id)
}
