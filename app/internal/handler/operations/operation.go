package operations

import (
	"encoding/json"
	"finance-manager-api-service/internal/apperror"
	"finance-manager-api-service/internal/client/operation_service/operation"
	h "finance-manager-api-service/internal/handler"
	"finance-manager-api-service/pkg/jwt"
	"finance-manager-api-service/pkg/logging"
	"finance-manager-api-service/pkg/utils"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const (
	operationsURL    = "/api/operations"
	operationByIdURL = "/api/operations/:uuid"
)

type operationHandler struct {
	Logger           *logging.Logger
	OperationService operation.Service
}

func NewOperationHandler(logger *logging.Logger, operationService operation.Service) h.Handler {
	return &operationHandler{
		Logger:           logger,
		OperationService: operationService,
	}
}

func (h *operationHandler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, operationsURL, jwt.Middleware(apperror.Middleware(h.CreateOperation)))
	router.HandlerFunc(http.MethodGet, operationByIdURL, jwt.Middleware(apperror.Middleware(h.GetOperationByUUID)))
	router.HandlerFunc(http.MethodPatch, operationByIdURL, jwt.Middleware(apperror.Middleware(h.PartiallyUpdateOperation)))
	router.HandlerFunc(http.MethodDelete, operationByIdURL, jwt.Middleware(apperror.Middleware(h.DeleteOperation)))
}

func (h *operationHandler) CreateOperation(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	defer utils.CloseBody(h.Logger, r.Body)

	var createdOperation operation.CreateOperationDTO
	if err := json.NewDecoder(r.Body).Decode(&createdOperation); err != nil {
		return apperror.BadRequestError("invalid JSON body")
	}

	operationUUID, err := h.OperationService.Create(r.Context(), createdOperation)
	if err != nil {
		return err
	}

	w.Header().Set("Location", fmt.Sprintf("%s/%s", operationsURL, operationUUID))
	w.WriteHeader(http.StatusCreated)
	return nil
}

func (h *operationHandler) GetOperationByUUID(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	operationUUID := params.ByName("uuid")

	op, err := h.OperationService.GetByUUID(r.Context(), operationUUID)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(op)
	return nil
}

func (h *operationHandler) PartiallyUpdateOperation(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	operationUUID := params.ByName("uuid")

	var updatedOperation operation.UpdateOperationDTO
	defer utils.CloseBody(h.Logger, r.Body)
	if err := json.NewDecoder(r.Body).Decode(&updatedOperation); err != nil {
		return apperror.BadRequestError("invalid JSON body")
	}

	if err := h.OperationService.Update(r.Context(), operationUUID, updatedOperation); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *operationHandler) DeleteOperation(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	operationUUID := params.ByName("uuid")

	if err := h.OperationService.Delete(r.Context(), operationUUID); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
