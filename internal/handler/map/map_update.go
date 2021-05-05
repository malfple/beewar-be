package mymap

import (
	"encoding/json"
	"gitlab.com/beewar/beewar-be/internal/controller/auth"
	"gitlab.com/beewar/beewar-be/internal/controller/mapmanager"
	"gitlab.com/beewar/beewar-be/internal/logger"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

// HandleMapUpdate handles map update
func HandleMapUpdate(w http.ResponseWriter, r *http.Request) {
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
	req := &MapUpdateRequest{}
	err = json.Unmarshal(body, req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp := &MapUpdateResponse{}

	err = mapmanager.UpdateMap(
		userID,
		req.MapID,
		req.MapType,
		req.Height,
		req.Width,
		req.Name,
		req.PlayerCount,
		req.TerrainInfo,
		req.UnitInfo)
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

// MapUpdateRequest is a request struct
type MapUpdateRequest struct {
	MapID       uint64 `json:"map_id"`
	MapType     uint8  `json:"map_type"`
	Height      int    `json:"height"`
	Width       int    `json:"width"`
	Name        string `json:"name"`
	PlayerCount uint8  `json:"player_count"`
	TerrainInfo []byte `json:"terrain_info"`
	UnitInfo    []byte `json:"unit_info"`
}

// MapUpdateResponse is a response struct
type MapUpdateResponse struct {
	ErrMsg string `json:"err_msg"`
}
