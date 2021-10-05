package campaign

import (
	"encoding/json"
	"gitlab.com/beewar/beewar-be/internal/controller/auth"
	"gitlab.com/beewar/beewar-be/internal/controller/campaign"
	"gitlab.com/beewar/beewar-be/internal/logger"
	"go.uber.org/zap"
	"net/http"
)

// HandleCurrentCampaign handles request to get current campaign
func HandleCurrentCampaign(w http.ResponseWriter, r *http.Request) {
	accessToken := r.Header.Get(auth.AccessTokenHeaderName)
	userID, _, err := auth.ValidateJWT(accessToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	gameID, err := campaign.GetCurrentCampaign(userID)

	resp := &CurrentCampaignResponse{
		GameID: gameID,
	}
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

// CurrentCampaignResponse is a response schema
type CurrentCampaignResponse struct {
	GameID uint64 `json:"game_id"`
	ErrMsg string `json:"err_msg"`
}
