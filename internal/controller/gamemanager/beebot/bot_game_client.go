package beebot

import (
	"fmt"
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager"
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager/message"
)

// BotGameClient is a modified client for internal bot use.
type BotGameClient struct {
	UserID     uint64
	Hub        *gamemanager.GameHub
	Replies    chan *message.GameMessage
	isShutDown bool
	shutdownC  chan bool
}

// NewBotGameClient creates a new bot client.
func NewBotGameClient(userID uint64) *BotGameClient {
	// does not initialize hub
	return &BotGameClient{
		UserID:     userID,
		Replies:    make(chan *message.GameMessage, 10),
		isShutDown: false,
		shutdownC:  make(chan bool),
	}
}

// GetUserID see function from gamemanager.GameClient
func (client *BotGameClient) GetUserID() uint64 {
	return client.UserID
}

// SetHub see function from gamemanager.GameClient
func (client *BotGameClient) SetHub(hub *gamemanager.GameHub) {
	client.Hub = hub
}

// GetHub see function from gamemanager.GameClient
func (client *BotGameClient) GetHub() *gamemanager.GameHub {
	return client.Hub
}

// SendMessageBack see function from gamemanager.GameClient
func (client *BotGameClient) SendMessageBack(msg *message.GameMessage) error {
	if !client.isShutDown {
		client.Replies <- msg
	}
	return nil
}

// Close see function from gamemanager.GameClient
func (client *BotGameClient) Close() {
	if !client.isShutDown {
		client.shutdownC <- true
	}
}

// Listen listens for incoming message and does the bot stuff
func (client *BotGameClient) Listen() {
	for !client.isShutDown {
		select {
		case <-client.shutdownC:
			client.isShutDown = true
			return
		case msg := <-client.Replies:
			fmt.Printf("beebot game id: %v, receive message cmd: %v\n", client.Hub.GameID, msg.Cmd)
		}
	}
}
