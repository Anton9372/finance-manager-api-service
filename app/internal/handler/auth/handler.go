package auth

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
	authURL   = "/api/auth"
	signUpURL = "/api/signup"
)

type handler struct {
	Logger      *logging.Logger
	UserService user_service.UserService
	JWTHelper   jwt.Helper
}

func NewAuthHandler(logger *logging.Logger, userService user_service.UserService, jwtHelper jwt.Helper) h.Handler {
	return &handler{
		Logger:      logger,
		UserService: userService,
		JWTHelper:   jwtHelper,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, authURL, apperror.Middleware(h.Auth))
	router.HandlerFunc(http.MethodPut, authURL, apperror.Middleware(h.Auth))
	router.HandlerFunc(http.MethodPost, signUpURL, apperror.Middleware(h.SignUp))
}

// SignUp
// @Summary 	Register user
// @Description Register user
// @Tags 		Auth
// @Accept		json
// @Produce 	json
// @Param 		input	body 	 user_service.SignUpUserDTO	true	"User's data"
// @Success 	201 	{object} jwt.TokenAndRefreshToken
// @Failure 	400 	{object} apperror.AppError "Validation error"
// @Failure 	418 	{object} apperror.AppError "Something wrong with application logic"
// @Failure 	500 	{object} apperror.AppError "Internal server error"
// @Router /signup [post]
func (h *handler) SignUp(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("Sign up")
	w.Header().Set("Content-Type", "application/json")
	defer utils.CloseBody(h.Logger, r.Body)

	var dto user_service.SignUpUserDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		return apperror.BadRequestError("invalid JSON body")
	}

	user, err := h.UserService.Create(r.Context(), dto)
	if err != nil {
		return err
	}

	token, err := h.JWTHelper.GenerateAccessToken(user)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(token)
	return nil
}

// Auth
// @Summary     Auth user and generate tokens or update refresh token
// @Description Auth user (POST) or update refresh token (PUT) and generate access token.
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Param       user         body       user_service.SignInUserDTO  false    "User's data"
// @Param       token        body       jwt.RefreshToken			false    "RefreshToken"
// @Success 	201 		{object} 	jwt.TokenAndRefreshToken
// @Failure     400         {object}    apperror.AppError   "Bad request or invalid JSON body"
// @Failure     401         {object}    apperror.AppError   "Unauthorized: invalid credentials"
// @Failure     500         {object}    apperror.AppError   "Internal server error"
// @Router      /auth       [post]
// @Router      /auth       [put]
func (h *handler) Auth(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("Auth")
	w.Header().Set("Content-Type", "application/json")
	defer utils.CloseBody(h.Logger, r.Body)

	var token []byte
	switch r.Method {
	case http.MethodPost:
		var dto user_service.SignInUserDTO
		if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
			return apperror.BadRequestError("invalid JSON body")
		}

		user, err := h.UserService.GetByEmailAndPassword(r.Context(), dto.Email, dto.Password)
		if err != nil {
			return err
		}

		token, err = h.JWTHelper.GenerateAccessToken(user)
		if err != nil {
			return err
		}
	case http.MethodPut:
		var rt jwt.RefreshToken

		err := json.NewDecoder(r.Body).Decode(&rt)
		if err != nil {
			return apperror.BadRequestError("failed to decode token")
		}

		token, err = h.JWTHelper.UpdateRefreshToken(rt)
		if err != nil {
			return err
		}
	}
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(token)
	return nil
}
