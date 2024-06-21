package jwt

import (
	"encoding/json"
	"finance-manager-api-service/internal/client/user_service"
	"finance-manager-api-service/internal/config"
	"finance-manager-api-service/pkg/cache"
	"finance-manager-api-service/pkg/logging"
	"github.com/cristalhq/jwt/v3"
	"github.com/google/uuid"
	"time"
)

type RefreshToken struct {
	RefreshToken string `json:"refresh_token"`
}

type Helper interface {
	GenerateAccessToken(u user_service.User) ([]byte, error)
	UpdateRefreshToken(rt RefreshToken) ([]byte, error)
}

type UserClaims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
}

type helper struct {
	RTCache cache.Repository
	logger  *logging.Logger
}

func NewHelper(rtCache cache.Repository, logger *logging.Logger) Helper {
	return &helper{RTCache: rtCache, logger: logger}
}

func (h helper) GenerateAccessToken(u user_service.User) ([]byte, error) {
	key := []byte(config.GetConfig().JWT.Secret)
	signer, err := jwt.NewSignerHS(jwt.HS256, key)
	if err != nil {
		return nil, err
	}
	builder := jwt.NewBuilder(signer)

	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        u.UUID,
			Audience:  []string{"users"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(60 * time.Minute)),
		},
		Email: u.Email,
	}

	token, err := builder.Build(claims)
	if err != nil {
		return nil, err
	}

	h.logger.Info("create token")
	refreshTokenUUID := uuid.New()
	userBytes, _ := json.Marshal(u)
	err = h.RTCache.Set([]byte(refreshTokenUUID.String()), userBytes, -1)
	if err != nil {
		h.logger.Error(err)
		return nil, err
	}

	jsonByte, err := json.Marshal(map[string]string{
		"token":         token.String(),
		"refresh_token": refreshTokenUUID.String(),
	})
	if err != nil {
		return nil, err
	}

	return jsonByte, nil
}

func (h helper) UpdateRefreshToken(rt RefreshToken) ([]byte, error) {
	defer h.RTCache.Del([]byte(rt.RefreshToken))

	userBytes, err := h.RTCache.Get([]byte(rt.RefreshToken))
	if err != nil {
		return nil, err
	}

	var u user_service.User
	err = json.Unmarshal(userBytes, &u)
	if err != nil {
		return nil, err
	}

	return h.GenerateAccessToken(u)
}
