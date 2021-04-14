package game

import (
	"encoding/json"
	"gitlab.com/beewar/beewar-be/internal/controller/auth"
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager"
	"gitlab.com/beewar/beewar-be/internal/logger"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

// HandleGameRegister handles game creation
func HandleGameRegister(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		logger.GetLogger().Error("error parse form", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	accessToken := r.Header.Get(auth.AccessTokenHeaderName)
	userID, _, err := auth.ValidateJWT(accessToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	gameID, err := strconv.ParseInt(r.Form.Get("game_id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	playerOrder, err := strconv.ParseInt(r.Form.Get("player_order"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp := &RegisterResponse{}

	err = gamemanager.RegisterForGame(userID, uint64(gameID), uint8(playerOrder))
	if err != nil {
		resp.ErrMsg = err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		logger.GetLogger().Error("error encode", zap.Error(err))
	}
}

// RegisterResponse is a response struct
type RegisterResponse struct {
	ErrMsg string `json:"err_msg"`
}
