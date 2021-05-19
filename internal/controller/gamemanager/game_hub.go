package gamemanager

import (
	"errors"
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager/loader"
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager/message"
	"gitlab.com/beewar/beewar-be/internal/logger"
	"go.uber.org/zap"
	"sync"
	"time"
)

var (
	// ErrClientDuplicate is returned when a duplicate userID for a hub tries to register
	ErrClientDuplicate = errors.New("only one user per game hub allowed")
)

const (
	// ErrMsgNotPlayer is returned when a non-player sends message to game hub
	ErrMsgNotPlayer = "you are not a player in this game"

	pingInterval = 20 * time.Second
)

// the hub is responsible for broadcasting messages to the clients

// GameHub is the central broadcasting service for a game
type GameHub struct {
	GameID     uint64
	GameLoader *loader.GameLoader
	Clients    map[uint64]GameClient
	MessageBus chan *message.GameMessage
	mutex      sync.Mutex // this lock is used to make sure broadcast and client register/unreg is synced
	isShutdown bool
	shutdownC  chan bool
}

// NewGameHub initializes a new game hub with the correct game id. Provide an onShutdown function that will be triggered
// when the game hub decides to shut down.
func NewGameHub(gameID uint64) (*GameHub, error) {
	gameLoader, err := loader.NewGameLoader(gameID)
	if err != nil {
		return nil, err
	}
	return &GameHub{
		GameID:     gameID,
		GameLoader: gameLoader,
		Clients:    make(map[uint64]GameClient),
		MessageBus: make(chan *message.GameMessage),
		mutex:      sync.Mutex{},
		isShutdown: false,
		shutdownC:  make(chan bool),
	}, nil
}

// RegisterClient registers the client to this hub
func (hub *GameHub) RegisterClient(client GameClient) error {
	var err error = nil
	hub.mutex.Lock()
	if !hub.isShutdown {
		if _, ok := hub.Clients[client.GetUserID()]; ok { // client found
			err = ErrClientDuplicate
		} else {
			hub.Clients[client.GetUserID()] = client
		}
	}
	hub.mutex.Unlock()
	return err
}

// UnregisterClient unregisters the client
func (hub *GameHub) UnregisterClient(client GameClient) {
	hub.mutex.Lock()
	if !hub.isShutdown {
		delete(hub.Clients, client.GetUserID())
	}
	hub.mutex.Unlock()
}

// sends message to a client.
// WARNING: hub.mutex should already be locked
func (hub *GameHub) sendMessageToClient(client GameClient, msg *message.GameMessage) {
	err := client.SendMessageBack(msg)
	if err != nil {
		client.Close()
		delete(hub.Clients, client.GetUserID())
	}
}

// handle message. returns (resp, isBroadcast?)
// WARNING: hub.mutex should already be locked
func (hub *GameHub) handleMessage(msg *message.GameMessage) (*message.GameMessage, bool) {
	if msg.Cmd == message.CmdGameData { // any user can get game data
		return hub.GameLoader.GameData(), false
	} else if msg.Cmd == message.CmdChat {
		if _, ok := hub.GameLoader.UserIDToPlayerMap[msg.Sender]; !ok { // non-player
			return message.GameErrorMessage(ErrMsgNotPlayer), false
		}
		return msg, true
	} else {
		return hub.GameLoader.HandleMessage(msg)
	}
}

// ListenAndBroadcast handles broadcasting.
// Pass in a waitgroup to wait for hub to shutdown before exiting application.
// Only run this once.
func (hub *GameHub) ListenAndBroadcast(wg *sync.WaitGroup) {
	defer wg.Done()

	pingTicker := time.NewTicker(pingInterval)
	defer pingTicker.Stop()

	for !hub.isShutdown {
		select {
		case <-hub.shutdownC: // YES! SHUTDOWN!!!
			return
		case <-pingTicker.C:
			// ping
			hub.mutex.Lock()
			for _, client := range hub.Clients {
				hub.sendMessageToClient(client, &message.GameMessage{Cmd: message.CmdPing})
			}
			hub.mutex.Unlock()
		case msg := <-hub.MessageBus:
			hub.mutex.Lock()
			logger.GetLogger().Debug("game hub: receive message", zap.String("cmd", msg.Cmd), zap.Uint64("sender", msg.Sender))
			// process message
			resp, isBroadcast := hub.handleMessage(msg)
			// broadcast
			if isBroadcast {
				for _, client := range hub.Clients {
					hub.sendMessageToClient(client, resp)
				}
			} else {
				hub.sendMessageToClient(hub.Clients[msg.Sender], resp)
			}
			hub.mutex.Unlock()
		}
	}
}

// Shutdown closes all client connection and stops the hub
func (hub *GameHub) Shutdown() {
	hub.mutex.Lock()
	if !hub.isShutdown {
		hub.isShutdown = true
		hub.GameLoader.SaveToDB()
		hub.shutdownC <- true // trigger shutdown for the listening goroutine
		for _, client := range hub.Clients {
			client.Close()
			delete(hub.Clients, client.GetUserID())
		}
	}
	hub.mutex.Unlock()
}
