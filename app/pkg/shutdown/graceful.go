package shutdown

import (
	"finance-manager-api-service/pkg/logging"
	"io"
	"os"
	"os/signal"
)

func Graceful(signals []os.Signal, closeItems ...io.Closer) {
	logger := logging.GetLogger()

	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, signals...)
	sig := <-sigChannel
	logger.Infof("Caught signal %s. Shutting down...", sig)

	for _, closer := range closeItems {
		if err := closer.Close(); err != nil {
			logger.Errorf("failed to close %v: %v", closer, err)
		}
	}
}
