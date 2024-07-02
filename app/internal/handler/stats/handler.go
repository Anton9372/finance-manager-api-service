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

// GetReport
// @Summary 	Get report about user's financial operations
// @Description Retrieves a list of operations with support for filtering and sorting.
// @Security	JWTAuth
// @Tags 		Stats
// @Produce 	json
// @Param 		user_uuid 	  path 	   string false  "User UUID"
// @Param 		category_name path 	   string false  "Category name (supports operators: substr)"
// @Param 		type	 	  path 	   string false  "Category type"
// @Param 		category_id   path 	   string false  "Category ID"
// @Param 		description   path 	   string false  "Description (supports operators: substr)"
// @Param 		money_sum 	  path 	   string false  "Money sum (supports operators: eq, neq, lt, lte, gt, gte, between)"
// @Param 		date_time     path 	   string false  "Date and time of operation (supports operators: eq, between; format: yyyy-mm-dd)"
// @Param 		sort_by 	  path 	   string false  "Field to sort by (money_sum, date_time, description)"
// @Param 		sort_order 	  path 	   string false  "Sort order (asc, desc)"
// @Success 	200 		  {object} stats_service.Report "Report"
// @Failure 	401 		   								"Unauthorized"
// @Failure 	400 		  {object} apperror.AppError 	"Validation error in filter or sort parameters"
// @Failure 	418 		  {object} apperror.AppError 	"Something wrong with application logic"
// @Failure 	500 		  {object} apperror.AppError 	"Internal server error"
// @Router /stats [get]
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
