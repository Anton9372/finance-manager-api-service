package rest

import (
	"finance-manager-api-service/pkg/logging"
	"finance-manager-api-service/pkg/utils"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type APIError struct {
	Message          string `json:"message,omitempty"`
	ErrorCode        string `json:"error_code,omitempty"`
	DeveloperMessage string `json:"developer_message,omitempty"`
}

func (e *APIError) ToString() string {
	return fmt.Sprintf("Err Code: %s, Err: %s, Developer Err: %s", e.ErrorCode, e.Message, e.DeveloperMessage)
}

type APIResponse struct {
	IsOk     bool
	response *http.Response
	Error    APIError
}

func (r *APIResponse) Body() io.ReadCloser {
	return r.response.Body
}

func (r *APIResponse) ReadBody() ([]byte, error) {
	defer utils.CloseBody(logging.GetLogger(), r.response.Body)
	return io.ReadAll(r.response.Body)
}

func (r *APIResponse) StatusCode() int {
	return r.response.StatusCode
}

func (r *APIResponse) Location() (*url.URL, error) {
	return r.response.Location()
}
