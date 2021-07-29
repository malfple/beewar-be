package campaign

import (
	"encoding/json"
	"gitlab.com/beewar/beewar-be/internal/controller/auth"
	"gitlab.com/beewar/beewar-be/internal/controller/campaign"
	"gitlab.com/beewar/beewar-be/internal/logger"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

// HandleStartCampaign handles start campaign
func HandleStartCampaign(w http.ResponseWriter, r *http.Request) {
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
	req := &StartCampaignRequest{}
	err = json.Unmarshal(body, req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp := &StartCampaignResponse{}

	gameID, err := campaign.StartNewCampaign(userID, req.CampaignLevel)
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

// StartCampaignRequest is a request schema
type StartCampaignRequest struct {
	CampaignLevel int `json:"campaign_level"`
}

// StartCampaignResponse is a response schema
type StartCampaignResponse struct {
	GameID uint64 `json:"game_id"`
	ErrMsg string `json:"err_msg"`
}
