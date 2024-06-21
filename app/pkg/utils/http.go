package utils

import (
	"finance-manager-api-service/pkg/logging"
	"io"
)

func CloseBody(logger *logging.Logger, body io.ReadCloser) {
	err := body.Close()
	if err != nil {
		logger.Fatalf("Error closing request body: %v", err)
	}
}
