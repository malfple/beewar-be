package game

import (
	"github.com/gorilla/websocket"
	"gitlab.com/beewar/beewar-be/configs"
	"gitlab.com/beewar/beewar-be/internal/controller/auth"
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager"
	"gitlab.com/beewar/beewar-be/internal/logger"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		for _, allowedOrigin := range configs.GetServerConfig().AllowedOrigins {
			if origin == allowedOrigin {
				return true
			}
		}
		logger.GetLogger().Debug("ws: origin not allowed", zap.String("origin", origin))
		return false
	},
	Subprotocols: []string{"game_room"},
}

// HandleGameWS handles websocket connection for a game
func HandleGameWS(w http.ResponseWriter, r *http.Request) {
	// requested protocols has to be in the form of:
	// ["game_room", <access_token>]  --  2 string values
	protocols := websocket.Subprotocols(r)
	if len(protocols) != 2 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	accessToken := protocols[1]
	userID, _, err := auth.ValidateJWT(accessToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// need to include game id in params
	gameID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// setup client-hub connection
	client := &gamemanager.GameClient{
		UserID: userID,
	}
	err = gamemanager.StartClientSession(client, uint64(gameID))
	if err != nil {
		logger.GetLogger().Debug("error start client session", zap.Error(err))
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	// token valid, game exists. now upgrade to websocket
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.GetLogger().Error("ws: error upgrade", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer func() {
		_ = c.Close()
	}()

	client.WS = c
	logger.GetLogger().Debug("client start listening", zap.Uint64("user_id", userID), zap.Int64("game_id", gameID))
	client.Listen()
	logger.GetLogger().Debug("client stop listening", zap.Uint64("user_id", userID), zap.Int64("game_id", gameID))

	gamemanager.EndClientSession(client)
}
