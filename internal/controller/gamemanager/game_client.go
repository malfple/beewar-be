package gamemanager

import (
	"github.com/gorilla/websocket"
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager/message"
)

// the client is responsible for sending messages to the hub
// it should not send any messages back (handled by the hub)

// GameClient is the client that will connect to game hub to receive broadcasts
type GameClient interface {
	// GetUserID returns the user id
	GetUserID() uint64
	// SetHub sets the hub
	SetHub(hub *GameHub)
	// GetHub returns the hub
	GetHub() *GameHub
	// SendMessageBack send message back through websocket
	SendMessageBack(msg *message.GameMessage) error
	// Close closes the websocket connection
	Close()
}

// DefaultGameClient is the default client
type DefaultGameClient struct {
	UserID uint64
	WS     *websocket.Conn // the websocket connection
	Hub    *GameHub
}

// GetUserID see function from GameClient
func (client *DefaultGameClient) GetUserID() uint64 {
	return client.UserID
}

// SetHub see function from GameClient
func (client *DefaultGameClient) SetHub(hub *GameHub) {
	client.Hub = hub
}

// GetHub see function from GameClient
func (client *DefaultGameClient) GetHub() *GameHub {
	return client.Hub
}

// SendMessageBack see function from GameClient
func (client *DefaultGameClient) SendMessageBack(msg *message.GameMessage) error {
	rawMsg, err := message.MarshalGameMessage(msg)
	if err != nil {
		panic("shouldn't have errored when marshaling")
	}
	return client.WS.WriteMessage(websocket.TextMessage, rawMsg)
}

// Close see function from GameClient
func (client *DefaultGameClient) Close() {
	client.WS.Close()
}

// Listen listens for incoming messages and sends them to the hub
// also registers itself to the hub on start and unreg on stop
func (client *DefaultGameClient) Listen() {
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
