package auth

import (
	"encoding/json"
	"gitlab.com/otqee/otqee-be/internal/auth"
	"gitlab.com/otqee/otqee-be/internal/logger"
	"go.uber.org/zap"
	"net/http"
)

// HandleToken handles access token request from refresh token
func HandleToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// take refresh token from cookie
	refreshToken := ""
	if refreshTokenCookie, err := r.Cookie("otqee-rtoken"); err == nil {
		refreshToken = refreshTokenCookie.Value
	}

	username := auth.ValidateRefreshToken(refreshToken)
	if username == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	accessToken := auth.GenerateJWT(username)

	resp := &TokenResponse{
		Token: accessToken,
	}
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		logger.GetLogger().Error("error encode", zap.Error(err))
	}
}

// TokenResponse is a response for token handler
type TokenResponse struct {
	Token string `json:"token"`
}
