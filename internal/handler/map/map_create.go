package mymap

import (
	"encoding/json"
	"gitlab.com/beewar/beewar-be/internal/controller/auth"
	"gitlab.com/beewar/beewar-be/internal/controller/mapmanager"
	"gitlab.com/beewar/beewar-be/internal/logger"
	"go.uber.org/zap"
	"net/http"
)

// HandleMapCreate handles map creation request
func HandleMapCreate(w http.ResponseWriter, r *http.Request) {
	accessToken := r.Header.Get(auth.AccessTokenHeaderName)
	userID, _, err := auth.ValidateJWT(accessToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	resp := &MapCreateResponse{}

	mapID, err := mapmanager.CreateEmptyMap(userID)
	if err != nil {
		resp.ErrMsg = err.Error()
	} else {
		resp.MapID = mapID
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		logger.GetLogger().Error("error encode", zap.Error(err))
	}
}

// MapCreateResponse is a response struct
type MapCreateResponse struct {
	MapID  uint64 `json:"map_id"`
	ErrMsg string `json:"err_msg"`
}
