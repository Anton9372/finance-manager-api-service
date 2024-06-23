package rest

import (
	"encoding/json"
	"errors"
	"finance-manager-api-service/pkg/logging"
	"finance-manager-api-service/pkg/utils"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"sync"
)

type BaseClient struct {
	BaseURL    string
	HTTPClient *http.Client
	Logger     *logging.Logger
	mu         sync.Mutex
}

func (c *BaseClient) SendRequest(req *http.Request) (*APIResponse, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.HTTPClient == nil {
		return nil, errors.New("no http client")
	}

	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	response, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request. error: %w", err)
	}

	apiResponse := APIResponse{
		IsOk:     true,
		response: response,
	}

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusBadRequest {
		defer utils.CloseBody(c.Logger, response.Body)
		apiResponse.IsOk = false

		var apiErr APIError
		err = json.NewDecoder(response.Body).Decode(&apiErr)
		if err != nil {
			return &apiResponse, fmt.Errorf("failed to parse apperror from response body: %w", err)
		}
		apiResponse.Error = apiErr
	}

	return &apiResponse, nil
}

func (c *BaseClient) BuildURL(resource string, filters []FilterOptions) (string, error) {
	var resultURL string
	parsedURL, err := url.ParseRequestURI(c.BaseURL)
	if err != nil {
		return resultURL, fmt.Errorf("failed to parse base URL. error: %w", err)
	}
	parsedURL.Path = path.Join(parsedURL.Path, resource)

	if len(filters) > 0 {
		q := parsedURL.Query()
		for _, fo := range filters {
			q.Set(fo.Field, fo.ToString())
		}
		parsedURL.RawQuery = q.Encode()
	}

	return parsedURL.String(), nil
}

func (c *BaseClient) Close() error {
	c.HTTPClient = nil
	return nil
}
