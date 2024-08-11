package users

import (
	"encoding/json"
	"finance-manager-api-service/internal/apperror"
	"finance-manager-api-service/internal/client/user_service"
	h "finance-manager-api-service/internal/handler"
	"finance-manager-api-service/pkg/jwt"
	"finance-manager-api-service/pkg/logging"
	"finance-manager-api-service/pkg/utils"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const (
	userProfileURL = "/api/user/profile"
)

type userHandler struct {
	Logger      *logging.Logger
	UserService user_service.UserService
}

func NewUserHandler(logger *logging.Logger, userService user_service.UserService) h.Handler {
	return &userHandler{
		Logger:      logger,
		UserService: userService,
	}
}

func (h *userHandler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPatch, userProfileURL, jwt.Middleware(apperror.Middleware(h.PartiallyUpdateUser)))
	router.HandlerFunc(http.MethodDelete, userProfileURL, jwt.Middleware(apperror.Middleware(h.DeleteUser)))
}

// PartiallyUpdateUser
// @Summary 	Update user
// @Description Update user's profile
// @Security	JWTAuth
// @Tags 		User
// @Accept		json
// @Param 		input 	body 	user_service.UpdateUserDTO  true  "User's data"
// @Success 	204
// @Failure 	400 	{object} apperror.AppError "Validation error"
// @Failure 	401 		   					   "Unauthorized"
// @Failure 	404 	{object} apperror.AppError "User is not found"
// @Failure 	418 	{object} apperror.AppError "Something wrong with application logic"
// @Failure 	500 	{object} apperror.AppError "Internal server error"
// @Router /user/profile [patch]
func (h *userHandler) PartiallyUpdateUser(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	if r.Context().Value("user_uuid") == nil {
		h.Logger.Error("no user_uuid in context")
		return apperror.UnauthorizedError("")
	}
	userUUID := r.Context().Value("user_uuid").(string)

	var updatedUser user_service.UpdateUserDTO
	defer utils.CloseBody(h.Logger, r.Body)
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		return apperror.BadRequestError("invalid JSON body")
	}

	updatedUser.UUID = userUUID
	err := h.UserService.Update(r.Context(), updatedUser)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

// DeleteUser
// @Summary 	Delete user
// @Description Delete user
// @Security	JWTAuth
// @Tags 		User
// @Param 		uuid 	path 	 string 	true  "User's uuid"
// @Success 	204
// @Failure 	401 		   						"Unauthorized"
// @Failure 	404 	{object} apperror.AppError "User is not found"
// @Failure 	418 	{object} apperror.AppError "Something wrong with application logic"
// @Failure 	500 	{object} apperror.AppError "Internal server error"
// @Router /user/profile [delete]
func (h *userHandler) DeleteUser(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	if r.Context().Value("user_uuid") == nil {
		h.Logger.Error("no user_uuid in context")
		return apperror.UnauthorizedError("")
	}
	userUUID := r.Context().Value("user_uuid").(string)

	err := h.UserService.Delete(r.Context(), userUUID)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
