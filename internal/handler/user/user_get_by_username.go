package user

import (
	"encoding/json"
	"gitlab.com/otqee/otqee-be/internal/access"
	"gitlab.com/otqee/otqee-be/internal/access/model"
	"gitlab.com/otqee/otqee-be/internal/logger"
	"go.uber.org/zap"
	"net/http"
)

// HandleUserGetByUsername handles profile query
func HandleUserGetByUsername(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	user := access.QueryUserByUsername(username)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := &GetResponse{User: nil}
	if user != nil {
		resp.User = user
		resp.User.Password = ""
	}

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		logger.GetLogger().Error("error encode", zap.Error(err))
	}
}

// GetResponse is a response for user get handler
type GetResponse struct {
	User *model.User `json:"user"`
}
