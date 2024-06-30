package stats

import (
	"finance-manager-api-service/internal/apperror"
	"finance-manager-api-service/internal/client/stats_service"
	h "finance-manager-api-service/internal/handler"
	"finance-manager-api-service/pkg/jwt"
	"finance-manager-api-service/pkg/logging"
	"finance-manager-api-service/pkg/rest"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
)

const (
	statsURL = "/api/stats"
)

type handler struct {
	Logger  *logging.Logger
	Service stats_service.Service
}

func NewHandler(logger *logging.Logger, service stats_service.Service) h.Handler {
	return &handler{
		Logger:  logger,
		Service: service,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, statsURL, jwt.Middleware(apperror.Middleware(h.GetReport)))
}

func (h *handler) GetReport(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	if r.Context().Value("user_uuid") == nil {
		h.Logger.Error("no user_uuid in context")
		return apperror.UnauthorizedError("")
	}
	userUUID := r.Context().Value("user_uuid").(string)

	params := r.URL.Query()

	var filters []rest.FilterOptions
	for key, values := range params {
		for _, value := range values {
			var operator string
			var vals []string

			if strings.Contains(value, ":") {
				parts := strings.SplitN(value, ":", 2)
				operator = parts[0]
				vals = strings.Split(parts[1], ",")
			} else {
				operator = ""
				vals = []string{value}
			}

			filters = append(filters, rest.FilterOptions{
				Field:    key,
				Operator: operator,
				Values:   vals,
			})
		}
	}

	filters = append(filters, rest.FilterOptions{
		Field:    "user_uuid",
		Operator: "",
		Values:   []string{userUUID},
	})

	report, err := h.Service.GetReport(r.Context(), userUUID, filters)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(report)
	return nil
}
