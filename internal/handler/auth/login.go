package auth

import (
	"encoding/json"
	"gitlab.com/otqee/otqee-be/internal/auth"
	"gitlab.com/otqee/otqee-be/internal/logger"
	"go.uber.org/zap"
	"net/http"
)

// HandleLogin handles user login
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		logger.GetLogger().Error("error parse form", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	username := r.Form.Get("username")
	password := r.Form.Get("password")
	token, statusCode := auth.Login(username, password)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	logger.GetLogger().Debug("login",
		zap.String("username", username),
		zap.Int("status_code", statusCode))

	if statusCode == http.StatusOK {
		resp := &LoginResponse{
			Token: token,
		}
		err := json.NewEncoder(w).Encode(resp)
		if err != nil {
			logger.GetLogger().Error("error encode", zap.Error(err))
		}
	}
}

// LoginResponse is a response for login handler
type LoginResponse struct {
	Token string `json:"token"`
}
