package game

import (
	"encoding/json"
	"gitlab.com/beewar/beewar-be/internal/access"
	"gitlab.com/beewar/beewar-be/internal/access/model"
	"gitlab.com/beewar/beewar-be/internal/controller/auth"
	"gitlab.com/beewar/beewar-be/internal/logger"
	"go.uber.org/zap"
	"net/http"
)

// HandleMyGames handles request to get a list of games from a user
func HandleMyGames(w http.ResponseWriter, r *http.Request) {
	accessToken := r.Header.Get(auth.AccessTokenHeaderName)
	userID, _, err := auth.ValidateJWT(accessToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	resp := &ListResponse{
		GameUsers: access.QueryGameUsersByUserID(userID),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		logger.GetLogger().Error("error encode", zap.Error(err))
	}
}

// ListResponse is response struct for game list handler
type ListResponse struct {
	GameUsers []*model.GameUser `json:"game_users"`
}
