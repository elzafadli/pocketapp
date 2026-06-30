package toronto

import (
	"net/http"

	"monitorapp/config"

	"github.com/Unleash/unleash-client-go/v3"
	"github.com/Unleash/unleash-client-go/v3/context"
)

//go:generate mockgen -destination=mocks/toggle.go -package=mocks -source=toggle.go
type ToggleService interface {
	NewContext(userId string) unleash.FeatureOption
	IsEnabled(name string, opts ...unleash.FeatureOption) bool
}

type Toggle struct {
	*unleash.Client
	Conf *config.Config `inject:"config"`
}

func (t *Toggle) Startup() error {
	client, err := unleash.NewClient(
		unleash.WithListener(&unleash.DebugListener{}),
		unleash.WithAppName(t.Conf.Toggle.AppName),
		unleash.WithUrl(t.Conf.Toggle.URL),
		unleash.WithCustomHeaders(http.Header{"Authorization": {t.Conf.Toggle.Token}}),
	)

	if err != nil {
		return err
	}

	t.Client = client

	return nil
}

func (f *Toggle) Shutdown() error { return f.Client.Close() }

func (f *Toggle) NewContext(userId string) unleash.FeatureOption {
	ctx := context.Context{
		UserId: userId,
	}

	return unleash.WithContext(ctx)
}

func (f *Toggle) IsEnabled(name string, opts ...unleash.FeatureOption) bool {
	return f.Client.IsEnabled(name, opts...)
}
