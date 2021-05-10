package gamemanager

import (
	"github.com/gorilla/websocket"
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager/message"
)

// the client is responsible for sending messages to the hub
// it should not send any messages back (handled by the hub)

// GameClient is the client that will connect to game hub to receive broadcasts
type GameClient struct {
	UserID uint64
	WS     *websocket.Conn // the websocket connection
	Hub    *GameHub
}

// Listen listens for incoming messages and sends them to the hub
// also registers itself to the hub on start and unreg on stop
func (client *GameClient) Listen() {
	for {
		_, rawMsg, err := client.WS.ReadMessage()
		if err != nil {
			break
		}
		msg, err := message.UnmarshalAndValidateGameMessage(rawMsg, client.UserID)
		if err != nil {
			break
		}
		client.Hub.MessageBus <- msg
	}
}

// NewGameClient creates a new client and connects it to the game hub
func NewGameClient(userID uint64, ws *websocket.Conn, hub *GameHub) *GameClient {
	client := &GameClient{
		UserID: userID,
		WS:     ws,
		Hub:    hub,
	}
	return client
}
