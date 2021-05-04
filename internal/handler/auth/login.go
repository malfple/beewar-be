package auth

import (
	"encoding/json"
	"gitlab.com/beewar/beewar-be/internal/controller/auth"
	"gitlab.com/beewar/beewar-be/internal/logger"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

// HandleLogin handles user login
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	req := &LoginRequest{}
	err = json.Unmarshal(body, req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	refreshToken, accessToken, statusCode := auth.Login(req.Username, req.Password)

	w.Header().Set("Content-Type", "application/json")

	logger.GetLogger().Debug("login",
		zap.String("username", req.Username),
		zap.Int("status_code", statusCode))

	if statusCode != http.StatusOK {
		w.WriteHeader(statusCode)
		return
	}

	// refresh token will be returned in the form of cookies,
	// and access token will be returned directly in the body

	http.SetCookie(w, &http.Cookie{
		Name:     auth.RefreshTokenCookieName,
		Value:    refreshToken,
		MaxAge:   864000, // 10 day cookie expiry. The expiry time for refresh token should be lower
		Path:     "/api/auth",
		HttpOnly: true,
	})
	w.WriteHeader(http.StatusOK)

	resp := &LoginResponse{
		Token: accessToken,
	}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		logger.GetLogger().Error("error encode", zap.Error(err))
	}
}

// LoginRequest is a request struct for login handler
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse is a response for login handler
type LoginResponse struct {
	Token string `json:"token"`
}
