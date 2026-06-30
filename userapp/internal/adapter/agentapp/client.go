package agentapp

import (
	"context"
	"errors"
	"fmt"
	"time"

	"userapp/config"
	"userapp/internal/domain/agent"
	"userapp/internal/domain/pocket"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	Conf  *config.Config `inject:"config"`
	Resty *resty.Client
}

type summarizeRequest struct {
	Data *pocket.PocketItem `json:"data"`
}

type summarizeResponse struct {
	Success  bool     `json:"success"`
	Summary  string   `json:"summary"`
	TodoList []string `json:"todoList"`
	Error    string   `json:"error,omitempty"`
}

func (c *Client) Startup() error {
	client := resty.New()
	client.SetBaseURL(c.Conf.Agent.Url).SetTimeout(30 * time.Second)
	
	c.Resty = client
	return nil
}

func (c *Client) Shutdown() error {
	return nil
}

func (c *Client) Summarize(ctx context.Context, item *pocket.PocketItem) (*agent.SummaryResponse, error) {
	reqBody := summarizeRequest{
		Data: item,
	}

	var resData summarizeResponse

	response, err := c.Resty.R().
		SetContext(ctx).
		SetBody(reqBody).
		SetResult(&resData).
		Post("/summarize")

	if err != nil {
		return nil, fmt.Errorf("failed to send summarize request: %w", err)
	}

	if response.IsError() {
		return nil, errors.New(response.String())
	}

	if !resData.Success {
		return nil, fmt.Errorf("agent API error: %s", resData.Error)
	}

	return &agent.SummaryResponse{
		Summary:  resData.Summary,
		TodoList: resData.TodoList,
	}, nil
}
