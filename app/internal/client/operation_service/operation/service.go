package operation

import (
	"bytes"
	"context"
	"encoding/json"
	"finance-manager-api-service/internal/apperror"
	"finance-manager-api-service/pkg/logging"
	"finance-manager-api-service/pkg/rest"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const requestWaitTime = 5 * time.Second

type Service interface {
	Create(ctx context.Context, dto CreateOperationDTO) (string, error)
	GetByUUID(ctx context.Context, uuid string) ([]byte, error)
	Update(ctx context.Context, uuid string, dto UpdateOperationDTO) error
	Delete(ctx context.Context, uuid string) error
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

func (c *client) Create(ctx context.Context, dto CreateOperationDTO) (string, error) {
	c.base.Logger.Info("Create operation")

	c.base.Logger.Debug("build url")
	url, err := c.base.BuildURL(c.Resource, nil)
	if err != nil {
		return "", fmt.Errorf("failed to build url: %w", err)
	}
	c.base.Logger.Tracef("url: %s", url)

	c.base.Logger.Debug("marshal dto to bytes")
	dataBytes, err := json.Marshal(dto)
	if err != nil {
		return "", fmt.Errorf("failed to marshal dto: %w", err)
	}

	c.base.Logger.Debug("create request")
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(dataBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	c.base.Logger.Debug("send request")
	reqCtx, cancel := context.WithTimeout(ctx, requestWaitTime)
	defer cancel()
	req = req.WithContext(reqCtx)
	response, err := c.base.SendRequest(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}

	if !response.IsOk {
		return "", apperror.APIError(response.Error.Code, response.Error.Message, response.Error.DeveloperMessage)
	}
	c.base.Logger.Debug("parse location header")
	operationURL, err := response.Location()
	if err != nil {
		return "", fmt.Errorf("failed to get location: %w", err)
	}
	c.base.Logger.Tracef("Location: %s", operationURL.String())

	splitURL := strings.Split(operationURL.String(), "/")
	operationUUID := splitURL[len(splitURL)-1]
	_, err = c.GetByUUID(ctx, operationUUID)
	return operationUUID, err
}

func (c *client) GetByUUID(ctx context.Context, uuid string) ([]byte, error) {
	c.base.Logger.Info("Get operation by uuid")
	var operation []byte

	c.base.Logger.Debug("build url")
	url, err := c.base.BuildURL(fmt.Sprintf("%s/%s", c.Resource+"/one", uuid), nil)
	if err != nil {
		return operation, fmt.Errorf("failed to build url: %w", err)
	}
	c.base.Logger.Tracef("url: %s", url)

	c.base.Logger.Debug("create request")
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return operation, fmt.Errorf("failed to create request: %w", err)
	}

	c.base.Logger.Debug("send request")
	reqCtx, cancel := context.WithTimeout(ctx, requestWaitTime)
	defer cancel()
	req = req.WithContext(reqCtx)
	response, err := c.base.SendRequest(req)
	if err != nil {
		return operation, fmt.Errorf("failed to send request: %w", err)
	}

	if !response.IsOk {
		return operation, apperror.APIError(response.Error.Code, response.Error.Message, response.Error.DeveloperMessage)
	}
	c.base.Logger.Debug("read response body")
	operation, err = response.ReadBody()
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}
	return operation, nil
}

func (c *client) Update(ctx context.Context, uuid string, dto UpdateOperationDTO) error {
	c.base.Logger.Info("Update operation")

	c.base.Logger.Debug("build url")
	url, err := c.base.BuildURL(fmt.Sprintf("%s/%s", c.Resource+"/one", uuid), nil)
	if err != nil {
		return fmt.Errorf("failed to build url: %w", err)
	}
	c.base.Logger.Tracef("url: %s", url)

	c.base.Logger.Debug("marshal dto to bytes")
	dataBytes, err := json.Marshal(dto)
	if err != nil {
		return fmt.Errorf("failed to marshal dto: %w", err)
	}

	c.base.Logger.Debug("create request")
	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(dataBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	c.base.Logger.Debug("send request")
	reqCtx, cancel := context.WithTimeout(ctx, requestWaitTime)
	defer cancel()
	req = req.WithContext(reqCtx)
	response, err := c.base.SendRequest(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	if !response.IsOk {
		return apperror.APIError(response.Error.Code, response.Error.Message, response.Error.DeveloperMessage)
	}
	return nil
}

func (c *client) Delete(ctx context.Context, uuid string) error {
	c.base.Logger.Info("Delete operation")

	c.base.Logger.Debug("build url")
	url, err := c.base.BuildURL(fmt.Sprintf("%s/%s", c.Resource+"/one", uuid), nil)
	if err != nil {
		return fmt.Errorf("failed to build url: %w", err)
	}
	c.base.Logger.Tracef("url: %s", url)

	c.base.Logger.Debug("create request")
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	c.base.Logger.Debug("send request")
	reqCtx, cancel := context.WithTimeout(ctx, requestWaitTime)
	defer cancel()
	req = req.WithContext(reqCtx)
	response, err := c.base.SendRequest(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	if !response.IsOk {
		return apperror.APIError(response.Error.Code, response.Error.Message, response.Error.DeveloperMessage)
	}
	return nil
}
