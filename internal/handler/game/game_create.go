package game

import (
	"encoding/json"
	"gitlab.com/beewar/beewar-be/internal/controller/auth"
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager"
	"gitlab.com/beewar/beewar-be/internal/logger"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

// HandleGameCreate handles game creation
func HandleGameCreate(w http.ResponseWriter, r *http.Request) {
	accessToken := r.Header.Get(auth.AccessTokenHeaderName)
	_, _, err := auth.ValidateJWT(accessToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	req := &CreateRequest{}
	err = json.Unmarshal(body, req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp := &CreateResponse{}

	gameID, err := gamemanager.CreateGame(req.MapID, req.Password)
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

// CreateRequest is a response struct
type CreateRequest struct {
	MapID    uint64 `json:"map_id"`
	Password string `json:"password"`
}

// CreateResponse is a response struct
type CreateResponse struct {
	GameID uint64 `json:"game_id"`
	ErrMsg string `json:"err_msg"`
}
