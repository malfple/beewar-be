package game

import (
	"github.com/gorilla/websocket"
	"gitlab.com/otqee/otqee-be/internal/access"
	"gitlab.com/otqee/otqee-be/internal/auth"
	"gitlab.com/otqee/otqee-be/internal/gamemanager"
	"gitlab.com/otqee/otqee-be/internal/logger"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// localhost route
		if r.Header.Get("Origin") == "http://localhost:3000" {
			return true
		}
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
	game := access.QueryGameByID(gameID)
	if game == nil {
		w.WriteHeader(http.StatusNoContent)
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

	client := gamemanager.NewGameClientByID(userID, c, gameID)
	client.Listen()
}
