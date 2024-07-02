package metric

import (
	h "finance-manager-api-service/internal/handler"
	"finance-manager-api-service/pkg/logging"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const (
	URL = "/api/heartbeat"
)

type handler struct {
	Logger *logging.Logger
}

func NewHandler(logger *logging.Logger) h.Handler {
	return &handler{
		Logger: logger,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, URL, h.Heartbeat)
}

// Heartbeat
// @Summary 	Heartbeat
// @Description Checks that the server is up and running
// @Tags 		Heartbeat
// @Success 	204
// @Router 		/metric [get]
func (h *handler) Heartbeat(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(204)
}
