package jwt

import (
	"context"
	"encoding/json"
	"finance-manager-api-service/internal/config"
	"finance-manager-api-service/pkg/logging"
	"github.com/cristalhq/jwt/v3"
	"net/http"
	"strings"
	"time"
)

func Middleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logging.GetLogger()
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")

		if len(authHeader) != 2 {
			logger.Error("Malformed token")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("malformed token"))
			return
		}

		logger.Debug("create jwt verifier")
		jwtToken := authHeader[1]
		key := []byte(config.GetConfig().JWT.Secret)
		verifier, err := jwt.NewVerifierHS(jwt.HS256, key)
		if err != nil {
			unauthorized(w, err)
			return
		}

		logger.Debug("parse and verify jwt token")
		token, err := jwt.ParseAndVerifyString(jwtToken, verifier)
		if err != nil {
			unauthorized(w, err)
			return
		}

		logger.Debug("parse user claims")
		var uc UserClaims
		err = json.Unmarshal(token.RawClaims(), &uc)
		if err != nil {
			unauthorized(w, err)
			return
		}

		if valid := uc.IsValidAt(time.Now()); !valid {
			logger.Error("token has been expired")
			unauthorized(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), "user_uuid", uc.ID)
		h(w, r.WithContext(ctx))
	}
}

func unauthorized(w http.ResponseWriter, err error) {
	logging.GetLogger().Error(err)
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte("unauthorized"))
}
