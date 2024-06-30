package stats_service

import (
	"context"
	"finance-manager-api-service/internal/apperror"
	"finance-manager-api-service/pkg/logging"
	"finance-manager-api-service/pkg/rest"
	"fmt"
	"net/http"
	"time"
)

const requestWaitTime = 5 * time.Second

type Service interface {
	GetReport(ctx context.Context, userUUID string, options []rest.FilterOptions) ([]byte, error)
}

type client struct {
	base     rest.BaseClient
	Resource string
}

func NewService(baseURL, resource string, logger *logging.Logger) Service {
	return &client{
		Resource: resource,
		base: rest.BaseClient{
			BaseURL: baseURL,
			HTTPClient: &http.Client{
				Timeout: 10 * time.Second,
			},
			Logger: logger,
		},
	}
}

func (c *client) GetReport(ctx context.Context, userUUID string, options []rest.FilterOptions) ([]byte, error) {
	c.base.Logger.Info("Get stats report")
	var report []byte

	c.base.Logger.Debug("build url")
	url, err := c.base.BuildURL(c.Resource, options)
	if err != nil {
		return report, fmt.Errorf("failed to build url: %w", err)
	}
	c.base.Logger.Tracef("url: %s", url)

	c.base.Logger.Debug("create request")
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return report, fmt.Errorf("failed to create request: %w", err)
	}

	c.base.Logger.Debug("send request")
	reqCtx, cancel := context.WithTimeout(ctx, requestWaitTime)
	defer cancel()
	req = req.WithContext(reqCtx)
	response, err := c.base.SendRequest(req)
	if err != nil {
		return report, fmt.Errorf("failed to send request: %w", err)
	}

	if !response.IsOk {
		appErr := apperror.APIError(response.Error.Code, response.Error.Message, response.Error.DeveloperMessage)
		appErr.WithParams(apperror.ErrorParams(response.Error.Params))
		appErr.WithFields(apperror.ErrorFields(response.Error.Fields))
		return report, appErr
	}
	c.base.Logger.Debug("read response body")
	report, err = response.ReadBody()
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}
	c.base.Logger.Debug("Get stats report successfully")
	return report, nil
}
