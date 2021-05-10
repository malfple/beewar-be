package gamemanager

import (
	"errors"
	"github.com/gorilla/websocket"
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager/loader"
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager/message"
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
	Clients    map[uint64]*GameClient
	MessageBus chan *message.GameMessage
	Mutex      sync.Mutex // this lock is used to make sure broadcast and client register/unreg is synced
	IsShutdown bool
	ShutdownC  chan bool
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
		Clients:    make(map[uint64]*GameClient),
		MessageBus: make(chan *message.GameMessage),
		Mutex:      sync.Mutex{},
		IsShutdown: false,
		ShutdownC:  make(chan bool),
	}, nil
}

// RegisterClient registers the client to this hub
func (hub *GameHub) RegisterClient(client *GameClient) error {
	var err error = nil
	hub.Mutex.Lock()
	if !hub.IsShutdown {
		if _, ok := hub.Clients[client.UserID]; ok { // client found
			err = ErrClientDuplicate
		} else {
			hub.Clients[client.UserID] = client
		}
	}
	hub.Mutex.Unlock()
	return err
}

// UnregisterClient unregisters the client
func (hub *GameHub) UnregisterClient(client *GameClient) {
	hub.Mutex.Lock()
	if !hub.IsShutdown {
		delete(hub.Clients, client.UserID)
	}
	hub.Mutex.Unlock()
}

// sends message to a client.
// WARNING: hub.Mutex should already be locked
func (hub *GameHub) sendMessageToClient(client *GameClient, msg *message.GameMessage) {
	rawMsg, err := message.MarshalGameMessage(msg)
	if err != nil {
		panic("shouldn't have errored when marshaling")
	}
	err = client.WS.WriteMessage(websocket.TextMessage, rawMsg)
	if err != nil {
		client.WS.Close()
		delete(hub.Clients, client.UserID)
	}
}

// handle message. returns (resp, isBroadcast?)
// WARNING: hub.Mutex should already be locked
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

	for !hub.IsShutdown {
		select {
		case <-hub.ShutdownC: // YES! SHUTDOWN!!!
			return
		case <-pingTicker.C:
			// ping
			hub.Mutex.Lock()
			for _, client := range hub.Clients {
				hub.sendMessageToClient(client, &message.GameMessage{Cmd: message.CmdPing})
			}
			hub.Mutex.Unlock()
		case msg := <-hub.MessageBus:
			hub.Mutex.Lock()
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
			hub.Mutex.Unlock()
		}
	}
}

// Shutdown closes all client connection and stops the hub
func (hub *GameHub) Shutdown() {
	hub.Mutex.Lock()
	if !hub.IsShutdown {
		hub.IsShutdown = true
		hub.GameLoader.SaveToDB()
		hub.ShutdownC <- true // trigger shutdown for the listening goroutine
		for _, client := range hub.Clients {
			client.WS.Close()
			delete(hub.Clients, client.UserID)
		}
	}
	hub.Mutex.Unlock()
}
