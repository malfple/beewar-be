package handler

import (
	"encoding/json"
	"gitlab.com/beewar/beewar-be/internal/auth"
	"gitlab.com/beewar/beewar-be/internal/gamemanager"
	"gitlab.com/beewar/beewar-be/internal/logger"
	"go.uber.org/zap"
	"net/http"
	"time"
)

var serverStartTime int64

func init() {
	serverStartTime = time.Now().Unix()
}

// HandleServerStats handles server stats query
func HandleServerStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := &ServerStatsResponse{
		HubCount:        gamemanager.GetHubCount(),
		SessionCount:    auth.GetRefreshTokenCount(),
		ServerStartTime: serverStartTime,
	}

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		logger.GetLogger().Error("error encode", zap.Error(err))
	}
}

// ServerStatsResponse is a response for server stats query
type ServerStatsResponse struct {
	HubCount        int   `json:"hub_count"`
	SessionCount    int   `json:"session_count"` // session here means refresh token count
	ServerStartTime int64 `json:"server_start_time"`
}
