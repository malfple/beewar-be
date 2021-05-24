package beebot

import (
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager"
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager/loader"
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager/message"
	"gitlab.com/beewar/beewar-be/internal/logger"
	"go.uber.org/zap"
	"time"
)

/*
This client works by sequence of messages

It will listen for messages, and only send messages on specific replies.
Therefore it needs a way to start itself when first initialized.

For now, it does a lot of illegal access to the hub's game loader, but whatever.
*/

// BotGameClient is a modified client for internal bot use.
type BotGameClient struct {
	UserID      uint64
	PlayerOrder int
	Hub         *gamemanager.GameHub
	Replies     chan *message.GameMessage
	isShutDown  bool
	shutdownC   chan bool
	nextTrigger string
}

// NewBotGameClient creates a new bot client.
func NewBotGameClient(userID uint64, playerOrder uint8) *BotGameClient {
	// does not initialize hub
	return &BotGameClient{
		UserID:      userID,
		PlayerOrder: int(playerOrder),
		Replies:     make(chan *message.GameMessage, 10),
		isShutDown:  false,
		shutdownC:   make(chan bool),
		nextTrigger: "",
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
	// pre-checks and kickstart if possible
	if client.isMyTurn() {
		client.doNextMove()
	}

	// listening loop
	for {
		select {
		case <-client.shutdownC:
			client.isShutDown = true
			return
		case msg := <-client.Replies:
			logger.GetLogger().Debug("beebot: receive message", zap.Uint64("game_id", client.Hub.GameID), zap.String("cmd", msg.Cmd))
			// handle replies
			if msg.Cmd == message.CmdPing {
				// periodically check if need to shutdown
				if client.isGameover() {
					client.isShutDown = true
					return
				}
			} else if client.nextTrigger == "" {
				// it's currently not my turn
				if msg.Cmd == message.CmdJoin || msg.Cmd == message.CmdEndTurn {
					// it could be my turn now!
					if client.isMyTurn() {
						client.doNextMove()
						break
					}
				}
			} else if msg.Cmd == client.nextTrigger {
				// reply received, proceed to next move
				client.doNextMove()
			}
		}
	}
}

// checks if beebot already lost, or game is over
func (client *BotGameClient) isGameover() bool {
	if client.Hub.GameLoader.Status == loader.GameStatusEnded {
		return true
	}
	return client.Hub.GameLoader.GameUsers[client.PlayerOrder-1].FinalTurns != 0 // if assigned final turns, gameover haha
}

// checks if it's currently beebot's turn to move
func (client *BotGameClient) isMyTurn() bool {
	// if game is already over, this will always return false
	if client.Hub.GameLoader.Status != loader.GameStatusOngoing {
		return false
	}

	return client.Hub.GameLoader.TurnPlayer == client.PlayerOrder
}

func (client *BotGameClient) sendMessageToHub(cmd string, data interface{}) {
	client.Hub.MessageBus <- &message.GameMessage{
		Cmd:    cmd,
		Sender: client.UserID,
		Data:   data,
	}
}

// the main function that sends actions to game hub
func (client *BotGameClient) doNextMove() {
	// safe delay before every move
	time.Sleep(time.Second)

	client.sendMessageToHub(message.CmdEndTurn, nil)

	client.nextTrigger = "" // if end turn, set the next trigger to empty string
}
