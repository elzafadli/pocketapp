package kaisel

import (
	"context"
	"errors"
	"fmt"
	"userapp/config"

	"github.com/go-resty/resty/v2"
	"github.com/runsystemid/golog"
)

var ErrorFailedToMigrate = errors.New("failed to migrate")

//go:generate mockgen -destination=mocks/kaisel.go -package=mocks -source=kaisel.go
type KaiselService interface {
	Migrate(request *MigrateRequest) error
}

type Kaisel struct {
	Conf  *config.Config `inject:"config"`
	Resty *resty.Client
}

type MigrateRequest struct {
	Schemas    []string `json:"schemas"`
	TenantName string   `json:"tenant_name"`
	TraceId    string   `json:"-"`
}

type MigrateRestyRequest struct {
	Endpoint       string
	MigrateRequest MigrateRequest
}

type MigrateResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (k *Kaisel) Startup() error {
	client := resty.New()
	client.SetBaseURL(k.Conf.Kaisel.URL).SetTimeout(k.Conf.Kaisel.Timeout)
	client.SetBasicAuth(k.Conf.Kaisel.User, k.Conf.Kaisel.Password)

	k.Resty = client
	return nil
}

func (r *Kaisel) Shutdown() error {
	return nil
}

func (k *Kaisel) Migrate(request *MigrateRequest) error {

	endpoints := []MigrateRestyRequest{
		{
			Endpoint:       "/v1/migrate",
			MigrateRequest: *request,
		},
		{
			Endpoint:       "/v1/seed",
			MigrateRequest: *request,
		},
	}

	for _, endpoint := range endpoints {

		var resp MigrateResponse

		response, err := k.Resty.R().
			SetBody(endpoint.MigrateRequest).
			SetResult(&resp).
			SetHeader("X-Correlation-ID", request.TraceId).
			Post(endpoint.Endpoint)

		if err != nil {
			return fmt.Errorf("failed to send request: %w", err)
		}

		if response.IsError() {
			return errors.New(response.String())
		}

		if resp.Status != "00" {
			golog.Error(context.Background(), "failed to migrate: %s", errors.New(resp.Message))
			return ErrorFailedToMigrate
		}
	}

	return nil
}
