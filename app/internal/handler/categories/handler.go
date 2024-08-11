package categories

import (
	"encoding/json"
	"finance-manager-api-service/internal/apperror"
	"finance-manager-api-service/internal/client/operation_service/category"
	h "finance-manager-api-service/internal/handler"
	"finance-manager-api-service/pkg/jwt"
	"finance-manager-api-service/pkg/logging"
	"finance-manager-api-service/pkg/utils"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const (
	categoriesURL   = "/api/categories"
	categoryByIdURL = "/api/categories/:uuid"
)

type categoryHandler struct {
	Logger          *logging.Logger
	CategoryService category.Service
}

func NewCategoryHandler(logger *logging.Logger, categoryService category.Service) h.Handler {
	return &categoryHandler{
		Logger:          logger,
		CategoryService: categoryService,
	}
}

func (h *categoryHandler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, categoriesURL, jwt.Middleware(apperror.Middleware(h.CreateCategory)))
	router.HandlerFunc(http.MethodGet, categoriesURL, jwt.Middleware(apperror.Middleware(h.GetCategories)))
	router.HandlerFunc(http.MethodPatch, categoryByIdURL, jwt.Middleware(apperror.Middleware(h.PartiallyUpdateCategory)))
	router.HandlerFunc(http.MethodDelete, categoryByIdURL, jwt.Middleware(apperror.Middleware(h.DeleteCategory)))
}

// CreateCategory
// @Summary 	Create category
// @Description Creates new category
// @Security	JWTAuth
// @Tags 		Category
// @Accept		json
// @Param 		input	body 	 category.CreateCategoryDTO	true	"Category data"
// @Success 	201
// @Failure 	401 		   						"Unauthorized"
// @Failure 	400 	{object} apperror.AppError "Validation error"
// @Failure 	418 	{object} apperror.AppError "Something wrong with application logic"
// @Failure 	500 	{object} apperror.AppError "Internal server error"
// @Router /categories [post]
func (h *categoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	if r.Context().Value("user_uuid") == nil {
		h.Logger.Error("no user_uuid in context")
		return apperror.UnauthorizedError("")
	}
	userUUID := r.Context().Value("user_uuid").(string)

	var createdCategory category.CreateCategoryDTO
	defer utils.CloseBody(h.Logger, r.Body)
	if err := json.NewDecoder(r.Body).Decode(&createdCategory); err != nil {
		return apperror.BadRequestError("invalid JSON body")
	}
	createdCategory.UserUUID = userUUID

	categoryUUID, err := h.CategoryService.Create(r.Context(), createdCategory)
	if err != nil {
		return err
	}

	w.Header().Set("Location", fmt.Sprintf("%s/%s", categoriesURL, categoryUUID))
	w.WriteHeader(http.StatusCreated)
	return nil
}

// GetCategories
// @Summary 	Get user's categories
// @Description Get list of categories belonging to user
// @Security	JWTAuth
// @Tags 		Category
// @Produce 	json
// @Param 		user_uuid 	path 	 string 	true   "User's uuid"
// @Success 	200			{object} []category.Category "Categories"
// @Failure 	401 		   						   "Unauthorized"
// @Failure 	404 		{object} apperror.AppError "User not found"
// @Failure 	418 		{object} apperror.AppError "Something wrong with application logic"
// @Failure 	500 		{object} apperror.AppError "Internal server error"
// @Router 		/categories	[get]
func (h *categoryHandler) GetCategories(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	if r.Context().Value("user_uuid") == nil {
		h.Logger.Error("no user_uuid in context")
		return apperror.UnauthorizedError("")
	}
	userUUID := r.Context().Value("user_uuid").(string)

	categories, err := h.CategoryService.GetByUserUUID(r.Context(), userUUID)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(categories)
	return nil
}

// PartiallyUpdateCategory
// @Summary 	Update category
// @Description Update category
// @Security	JWTAuth
// @Tags 		Category
// @Accept		json
// @Param 		uuid 		path 	 string 					true  "Category's uuid"
// @Param 		input 		body 	 category.UpdateCategoryDTO true  "Category's data"
// @Success 	204
// @Failure 	400 	{object} apperror.AppError "Validation error"
// @Failure 	401 		   					   "Unauthorized"
// @Failure 	404 	{object} apperror.AppError "Category is not found"
// @Failure 	418 	{object} apperror.AppError "Something wrong with application logic"
// @Failure 	500 	{object} apperror.AppError "Internal server error"
// @Router /categories/:uuid [patch]
func (h *categoryHandler) PartiallyUpdateCategory(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	var updatedCategory category.UpdateCategoryDTO
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	categoryUUID := params.ByName("uuid")
	defer utils.CloseBody(h.Logger, r.Body)
	if err := json.NewDecoder(r.Body).Decode(&updatedCategory); err != nil {
		return apperror.BadRequestError("invalid JSON body")
	}

	updatedCategory.UUID = categoryUUID
	err := h.CategoryService.Update(r.Context(), updatedCategory)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

// DeleteCategory
// @Summary 	Delete category
// @Description Delete category
// @Security	JWTAuth
// @Tags 		Category
// @Param 		uuid 	path 	 string 	true  "Category's uuid"
// @Success 	204
// @Failure 	401 		   						"Unauthorized"
// @Failure 	404 	{object} apperror.AppError "Category is not found"
// @Failure 	418 	{object} apperror.AppError "Something wrong with application logic"
// @Failure 	500 	{object} apperror.AppError "Internal server error"
// @Router /categories/:uuid [delete]
func (h *categoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	categoryUUID := params.ByName("uuid")

	err := h.CategoryService.Delete(r.Context(), categoryUUID)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
