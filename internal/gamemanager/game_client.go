package gamemanager

import (
	"github.com/gorilla/websocket"
	"gitlab.com/otqee/otqee-be/internal/gamemanager/message"
)

// the client is responsible for sending messages to the hub
// it should not send any messages back (handled by the hub)

// GameClient is the client that will connect to game hub to receive broadcasts
type GameClient struct {
	WS  *websocket.Conn // the websocket connection
	Hub *GameHub
}

// Listen listens for incoming messages and sends them to the hub
// also registers itself to the hub on start and unreg on stop
func (client *GameClient) Listen() {
	client.Hub.RegisterClient(client)
	for {
		_, rawMsg, err := client.WS.ReadMessage()
		if err != nil {
			break
		}
		msg, err := message.UnmarshalGameMessage(rawMsg)
		if err != nil {
			break
		}
		if msg.Cmd == "SHUTDOWN" { // well. lmao
			break
		}
		client.Hub.MessageBus <- msg
	}
	client.Hub.UnregisterClient(client)
}

// NewGameClient creates a new client and connects it to the game hub
func NewGameClient(ws *websocket.Conn, hub *GameHub) *GameClient {
	client := &GameClient{
		WS:  ws,
		Hub: hub,
	}
	return client
}
