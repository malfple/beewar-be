package handler

import (
	"encoding/json"
	"gitlab.com/otqee/otqee-be/internal/auth"
	"gitlab.com/otqee/otqee-be/internal/gamemanager"
	"gitlab.com/otqee/otqee-be/internal/logger"
	"go.uber.org/zap"
	"net/http"
)

// HandleServerStats handles server stats query
func HandleServerStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := &ServerStatsResponse{
		HubCount:     gamemanager.GetHubCount(),
		SessionCount: auth.GetRefreshTokenCount(),
	}

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		logger.GetLogger().Error("error encode", zap.Error(err))
	}
}

// ServerStatsResponse is a response for server stats query
type ServerStatsResponse struct {
	HubCount     int `json:"hub_count"`
	SessionCount int `json:"session_count"` // session here means refresh token count
}
