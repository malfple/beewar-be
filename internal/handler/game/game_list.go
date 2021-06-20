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

// HandleGameList handles request to get open/waiting games
func HandleGameList(w http.ResponseWriter, r *http.Request) {
	accessToken := r.Header.Get(auth.AccessTokenHeaderName)
	_, _, err := auth.ValidateJWT(accessToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	games, err := access.QueryWaitingGames()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp := &ListResponse{
		Games: games,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		logger.GetLogger().Error("error encode", zap.Error(err))
	}
}

// ListResponse is a response for game list handler
type ListResponse struct {
	Games []*model.Game `json:"games"`
}
