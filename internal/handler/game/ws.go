package game

import (
	"fmt"
	"github.com/gorilla/websocket"
	"gitlab.com/otqee/otqee-be/internal/logger"
	"go.uber.org/zap"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// localhost route
		if r.Header.Get("Origin") == "http://localhost:3000" {
			return true
		}
		return false
	},
}

// HandleGameWS handles websocket connection for a game
func HandleGameWS(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.GetLogger().Error("ws: error upgrade", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer func() {
		fmt.Println("ws closed")
		_ = c.Close()
	}()

	fmt.Println("ws opened")
	err = c.WriteMessage(websocket.TextMessage, []byte("server hello"))
	if err != nil {
		fmt.Println("error write: ", err)
	}
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			fmt.Println("error read: ", err)
			break
		}
		fmt.Printf("mt %v message %v\n", mt, message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			fmt.Println("error write: ", err)
		}
	}
}
