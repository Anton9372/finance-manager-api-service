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

// CreateOperation
// @Summary 	Create operation
// @Description Creates new operation
// @Security	JWTAuth
// @Tags 		Operation
// @Accept		json
// @Param 		input	body 	 operation.CreateOperationDTO	true	"Operation's data"
// @Success 	201
// @Failure 	401 		   						"Unauthorized"
// @Failure 	400 	{object} apperror.AppError "Validation error"
// @Failure 	418 	{object} apperror.AppError "Something wrong with application logic"
// @Failure 	500 	{object} apperror.AppError "Internal server error"
// @Router /operations [post]
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

// GetOperationByUUID
// @Summary 	Get operation by uuid
// @Description Get operation by uuid
// @Security	JWTAuth
// @Tags 		Operation
// @Produce 	json
// @Param 		uuid 	path 	 string 	true   		"Operation's uuid"
// @Success 	200		{object} operation.Operation  	"Operation"
// @Failure 	401 		   						"Unauthorized"
// @Failure 	404 	{object} apperror.AppError "Operation not found"
// @Failure 	418 	{object} apperror.AppError "Something wrong with application logic"
// @Failure 	500 	{object} apperror.AppError "Internal server error"
// @Router 		/operations/:uuid	[get]
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

// PartiallyUpdateOperation
// @Summary 	Update Operation
// @Description Update Operation
// @Security	JWTAuth
// @Tags 		Operation
// @Accept		json
// @Param 		uuid 		path 	 string 						true  "Operation's uuid"
// @Param 		input 		body 	 operation.UpdateOperationDTO 	true  "Operation's data"
// @Success 	204
// @Failure 	401 		   						"Unauthorized"
// @Failure 	400 	{object} apperror.AppError "Validation error"
// @Failure 	418 	{object} apperror.AppError "Something wrong with application logic"
// @Failure 	500 	{object} apperror.AppError "Internal server error"
// @Router /operations/:uuid [patch]
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

// DeleteOperation
// @Summary 	Delete operation
// @Description Delete operation
// @Security	JWTAuth
// @Tags 		Operation
// @Param 		uuid 	path 	 string 	true  "Operation's uuid"
// @Success 	204
// @Failure 	401 		   						"Unauthorized"
// @Failure 	404 	{object} apperror.AppError "Operation is not found"
// @Failure 	418 	{object} apperror.AppError "Something wrong with application logic"
// @Failure 	500 	{object} apperror.AppError "Internal server error"
// @Router /operations/:uuid [delete]
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
