package category

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
	Create(ctx context.Context, dto CreateCategoryDTO) (string, error)
	GetByUUID(ctx context.Context, uuid string) ([]byte, error)
	GetByUserUUID(ctx context.Context, userUUID string) ([]byte, error)
	Update(ctx context.Context, uuid string, dto UpdateCategoryDTO) error
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

func (c *client) Create(ctx context.Context, dto CreateCategoryDTO) (string, error) {
	c.base.Logger.Info("Create category")

	c.base.Logger.Debug("build url")
	url, err := c.base.BuildURL(c.Resource, nil)
	if err != nil {
		return "", fmt.Errorf("failed to build url: %w", err)
	}
	c.base.Logger.Tracef("url: %s", url)

	c.base.Logger.Debug("marshal dto to bytes")
	dataBytes, err := json.Marshal(dto)
	if err != nil {
		return "", fmt.Errorf("failed to marshal dto to bytes: %w", err)
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
	categoryURL, err := response.Location()
	if err != nil {
		return "", fmt.Errorf("failed to get location: %w", err)
	}
	c.base.Logger.Tracef("Location: %s", categoryURL.String())

	splitURL := strings.Split(categoryURL.String(), "/")
	categoryUUID := splitURL[len(splitURL)-1]
	_, err = c.GetByUUID(ctx, categoryUUID)
	return categoryUUID, err
}

func (c *client) GetByUUID(ctx context.Context, uuid string) ([]byte, error) {
	c.base.Logger.Info("Get category by uuid")
	var category []byte

	c.base.Logger.Debug("build url")
	url, err := c.base.BuildURL(fmt.Sprintf("%s/%s", c.Resource+"/one", uuid), nil)
	if err != nil {
		return category, fmt.Errorf("failed to builf url: %w", err)
	}
	c.base.Logger.Tracef("url: %s", url)

	c.base.Logger.Debug("create request")
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return category, fmt.Errorf("failed to create request: %w", err)
	}

	c.base.Logger.Debug("send request")
	reqCtx, cancel := context.WithTimeout(ctx, requestWaitTime)
	defer cancel()
	req = req.WithContext(reqCtx)
	response, err := c.base.SendRequest(req)
	if err != nil {
		return category, fmt.Errorf("failed to send request: %w", err)
	}

	if !response.IsOk {
		return category, apperror.APIError(response.Error.Code, response.Error.Message, response.Error.DeveloperMessage)
	}
	c.base.Logger.Debug("read response body")
	category, err = response.ReadBody()
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}
	return category, nil
}

func (c *client) GetByUserUUID(ctx context.Context, userUUID string) ([]byte, error) {
	c.base.Logger.Info("Get categories by user uuid")
	var categories []byte

	c.base.Logger.Debug("build url")
	url, err := c.base.BuildURL(fmt.Sprintf("%s/%s", c.Resource+"/user_uuid", userUUID), nil)
	if err != nil {
		return categories, fmt.Errorf("failed to build url: %w", err)
	}
	c.base.Logger.Tracef("url: %s", url)

	c.base.Logger.Debug("create request")
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return categories, fmt.Errorf("failed to create request: %w", err)
	}

	c.base.Logger.Debug("send request")
	reqCtx, cancel := context.WithTimeout(ctx, requestWaitTime)
	defer cancel()
	req = req.WithContext(reqCtx)
	response, err := c.base.SendRequest(req)
	if err != nil {
		return categories, fmt.Errorf("failed to send request: %w", err)
	}

	if !response.IsOk {
		return categories, apperror.APIError(response.Error.Code, response.Error.Message, response.Error.DeveloperMessage)
	}
	c.base.Logger.Debug("read response body")
	categories, err = response.ReadBody()
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}
	return categories, nil
}

func (c *client) Update(ctx context.Context, uuid string, dto UpdateCategoryDTO) error {
	c.base.Logger.Info("Update category")

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
		return fmt.Errorf("failed to create request; %w", err)
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
	c.base.Logger.Info("Delete category")

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
