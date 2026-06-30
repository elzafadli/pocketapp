package cache

import (
	"context"
	"fmt"
	"time"

	"monitorapp/internal/domain/shared/identity"
	"monitorapp/internal/domain/template"

	"github.com/runsystemid/gocache"
)

var keyTemplate = "template-%s"

type TemplateRepository struct {
	Cache gocache.Service     `inject:"cache"`
	Next  template.Repository `inject:"templateRepository"`
}

func (t *TemplateRepository) Create(ctx context.Context, data *template.Template) (*template.Template, error) {
	return t.Next.Create(ctx, data)
}

func (t *TemplateRepository) Update(ctx context.Context, data *template.Template) error {
	key := fmt.Sprintf(keyTemplate, data.ID.String())
	t.Cache.Delete(ctx, key)

	return t.Next.Update(ctx, data)
}

func (t *TemplateRepository) GetByID(ctx context.Context, id identity.ID) (*template.Template, error) {
	key := fmt.Sprintf(keyTemplate, id.String())
	data := &template.Template{}

	err := t.Cache.Get(ctx, key, &data)
	if err == nil {
		return data, nil
	}

	if err != gocache.ErrNil {
		return nil, err
	}

	data, err = t.Next.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	t.Cache.Put(ctx, key, data, 1*time.Minute)

	return data, nil
}
