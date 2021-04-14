package game

import (
	"encoding/json"
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager"
	"gitlab.com/beewar/beewar-be/internal/logger"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

// HandleGameCreate handles game creation
func HandleGameCreate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		logger.GetLogger().Error("error parse form", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	mapID, err := strconv.ParseInt(r.Form.Get("map_id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	password := r.Form.Get("password")

	resp := &CreateResponse{}

	gameID, err := gamemanager.CreateGame(uint64(mapID), password)
	if err != nil {
		resp.ErrMsg = err.Error()
	} else {
		resp.GameID = gameID
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		logger.GetLogger().Error("error encode", zap.Error(err))
	}
}

// CreateResponse is a response struct
type CreateResponse struct {
	GameID uint64 `json:"game_id"`
	ErrMsg string `json:"err_msg"`
}
