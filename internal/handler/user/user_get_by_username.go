package user

import (
	"encoding/json"
	"gitlab.com/beewar/beewar-be/internal/access/model"
	userController "gitlab.com/beewar/beewar-be/internal/controller/user"
	"gitlab.com/beewar/beewar-be/internal/logger"
	"go.uber.org/zap"
	"net/http"
)

// HandleUserGetByUsername handles profile query
func HandleUserGetByUsername(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := &GetResponse{
		User: userController.GetByUsername(r.URL.Query().Get("username")),
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
