package user_service

import (
	"bytes"
	"context"
	"encoding/json"
	"finance-manager-api-service/internal/apperror"
	"finance-manager-api-service/pkg/logging"
	"finance-manager-api-service/pkg/rest"
	"finance-manager-api-service/pkg/utils"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const requestWaitTime = 5 * time.Second

type UserService interface {
	Create(ctx context.Context, dto SignUpUserDTO) (User, error)
	GetByUUID(ctx context.Context, uuid string) (User, error)
	GetByEmailAndPassword(ctx context.Context, email, password string) (User, error)
	Update(ctx context.Context, uuid string, dto UpdateUserDTO) error
	Delete(ctx context.Context, uuid string) error
}

type client struct {
	base     rest.BaseClient
	Resource string
}

func NewService(baseURL string, resource string, logger *logging.Logger) UserService {
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

func (c *client) Create(ctx context.Context, dto SignUpUserDTO) (User, error) {
	c.base.Logger.Info("Create user")
	var user User

	c.base.Logger.Debug("build url")
	url, err := c.base.BuildURL(c.Resource, nil)
	if err != nil {
		return user, fmt.Errorf("failed to build url: %w", err)
	}
	c.base.Logger.Tracef("url: %s", url)

	c.base.Logger.Debug("marshal dto to bytes")
	dataBytes, err := json.Marshal(dto)
	if err != nil {
		return user, fmt.Errorf("failed to marshal dto: %w", err)
	}

	c.base.Logger.Debug("create request")
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(dataBytes))
	if err != nil {
		return user, fmt.Errorf("failed to create request: %w", err)
	}

	c.base.Logger.Debug("send request")
	reqCtx, cancel := context.WithTimeout(ctx, requestWaitTime)
	defer cancel()
	req = req.WithContext(reqCtx)
	response, err := c.base.SendRequest(req)
	if err != nil {
		return user, fmt.Errorf("failed to send request: %w", err)
	}

	if !response.IsOk {
		return user, apperror.APIError(response.Error.Code, response.Error.Message, response.Error.DeveloperMessage)
	}
	c.base.Logger.Debug("parse location header")
	userURL, err := response.Location()
	if err != nil {
		return user, fmt.Errorf("failed to get location: %w", err)
	}
	c.base.Logger.Tracef("Location: %s", userURL.String())

	splitURL := strings.Split(userURL.String(), "/")
	userUUID := splitURL[len(splitURL)-1]
	user, err = c.GetByUUID(ctx, userUUID)
	return user, err
}

func (c *client) GetByUUID(ctx context.Context, uuid string) (User, error) {
	c.base.Logger.Info("Get user by uuid")
	var user User

	c.base.Logger.Debug("build url")
	url, err := c.base.BuildURL(fmt.Sprintf("%s/%s", c.Resource+"/one", uuid), nil)
	if err != nil {
		return user, fmt.Errorf("failed to build url: %w", err)
	}
	c.base.Logger.Tracef("url: %s", url)

	c.base.Logger.Debug("create request")
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return user, fmt.Errorf("failed to create request: %w", err)
	}

	c.base.Logger.Debug("send request")
	reqCtx, cancel := context.WithTimeout(ctx, requestWaitTime)
	defer cancel()
	req = req.WithContext(reqCtx)
	response, err := c.base.SendRequest(req)
	if err != nil {
		return user, fmt.Errorf("failed to send request: %w", err)
	}

	if !response.IsOk {
		return user, apperror.APIError(response.Error.Code, response.Error.Message, response.Error.DeveloperMessage)
	}
	defer utils.CloseBody(c.base.Logger, response.Body())
	if err = json.NewDecoder(response.Body()).Decode(&user); err != nil {
		return user, fmt.Errorf("failed to decode response: %w", err)
	}
	return user, nil
}

func (c *client) GetByEmailAndPassword(ctx context.Context, email, password string) (User, error) {
	c.base.Logger.Info("Get user by email and password")
	var user User

	filters := []rest.FilterOptions{
		{
			Field:  "email",
			Values: []string{email},
		},
		{
			Field:  "password",
			Values: []string{password},
		},
	}

	c.base.Logger.Debug("build url")
	url, err := c.base.BuildURL(c.Resource, filters)
	if err != nil {
		return user, fmt.Errorf("failed to build url: %w", err)
	}
	c.base.Logger.Tracef("url: %s", url)

	c.base.Logger.Debug("create request")
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return user, fmt.Errorf("failed to create request: %w", err)
	}

	c.base.Logger.Debug("send request")
	reqCtx, cancel := context.WithTimeout(ctx, requestWaitTime)
	defer cancel()
	req = req.WithContext(reqCtx)
	response, err := c.base.SendRequest(req)
	if err != nil {
		return user, fmt.Errorf("failed to send request: %w", err)
	}

	if !response.IsOk {
		return user, apperror.APIError(response.Error.Code, response.Error.Message, response.Error.DeveloperMessage)
	}
	defer utils.CloseBody(c.base.Logger, response.Body())
	if err = json.NewDecoder(response.Body()).Decode(&user); err != nil {
		return user, fmt.Errorf("failed to decode response: %w", err)
	}
	return user, nil
}

func (c *client) Update(ctx context.Context, uuid string, dto UpdateUserDTO) error {
	//TODO
	c.base.Logger.Fatal("update user is not implemented yet")
	return fmt.Errorf("update user is not implemented yet")
}

func (c *client) Delete(ctx context.Context, uuid string) error {
	//TODO
	c.base.Logger.Fatal("delete user is not implemented yet")
	return fmt.Errorf("delete user is not implemented yet")
}
