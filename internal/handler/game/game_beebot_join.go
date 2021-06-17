package game

import (
	"encoding/json"
	"gitlab.com/beewar/beewar-be/internal/controller/auth"
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager/beebot"
	"gitlab.com/beewar/beewar-be/internal/logger"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

// HandleGameBeebotJoin handles beebot join game request
func HandleGameBeebotJoin(w http.ResponseWriter, r *http.Request) {
	accessToken := r.Header.Get(auth.AccessTokenHeaderName)
	userID, _, err := auth.ValidateJWT(accessToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	req := &BeebotJoinRequest{}
	err = json.Unmarshal(body, req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp := &BeebotJoinResponse{}

	resp.ErrMsg = beebot.AskBeebotToJoinGame(userID, req.GameID, req.PlayerOrder, req.Password)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		logger.GetLogger().Error("error encode", zap.Error(err))
	}
}

// BeebotJoinRequest is a request struct
type BeebotJoinRequest struct {
	GameID      uint64 `json:"game_id"`
	PlayerOrder uint8  `json:"player_order"`
	Password    string `json:"password"`
}

// BeebotJoinResponse is a response struct
type BeebotJoinResponse struct {
	ErrMsg string `json:"err_msg"`
}
