package user

import (
	"encoding/json"
	"gitlab.com/beewar/beewar-be/internal/access"
	"gitlab.com/beewar/beewar-be/internal/access/model"
	"gitlab.com/beewar/beewar-be/internal/logger"
	"go.uber.org/zap"
	"net/http"
)

// HandleUserGetByUsername handles profile query
func HandleUserGetByUsername(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	user, err := access.QueryUserByUsername(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := &GetByUsernameResponse{User: nil}
	if user != nil {
		resp.User = user
		resp.User.Password = "it's a secret haha"
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		logger.GetLogger().Error("error encode", zap.Error(err))
	}
}

// GetByUsernameResponse is a response for user get handler
type GetByUsernameResponse struct {
	User *model.User `json:"user"`
}
