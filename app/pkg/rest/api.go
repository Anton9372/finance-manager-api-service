package rest

import (
	"finance-manager-api-service/pkg/logging"
	"finance-manager-api-service/pkg/utils"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type ErrorFields map[string]string
type ErrorParams map[string]string

type APIError struct {
	Code             string      `json:"code,omitempty"`
	Message          string      `json:"message,omitempty"`
	DeveloperMessage string      `json:"developer_message,omitempty"`
	Fields           ErrorFields `json:"fields,omitempty"`
	Params           ErrorParams `json:"params,omitempty"`
}

func (e *APIError) ToString() string {
	return fmt.Sprintf("Err Code: %s, Err: %s, Developer Err: %s", e.Code, e.Message, e.DeveloperMessage)
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
