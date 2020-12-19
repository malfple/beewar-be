package gamemanager

import (
	"errors"
	"github.com/gorilla/websocket"
	"gitlab.com/otqee/otqee-be/internal/gamemanager/loader"
	"gitlab.com/otqee/otqee-be/internal/gamemanager/message"
	"sync"
)

var (
	// ErrClientDuplicate is returned when a duplicate userID for a hub tries to register
	ErrClientDuplicate = errors.New("only one user per game hub allowed")
)

const (
	// ErrMsgNotPlayer is returned when a non-player sends message to game hub
	ErrMsgNotPlayer = "you are not a player in this game"
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
	OnShutdown func() // called when not forced to shutdown
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
			// send game data
			gameDataMsg := hub.GameLoader.GameData()
			rawMsg, err := message.MarshalGameMessage(gameDataMsg)
			if err != nil {
				panic("shouldn't have errored when marshaling")
			}
			err = client.WS.WriteMessage(websocket.TextMessage, rawMsg)
			if err != nil {
				client.WS.Close()
				delete(hub.Clients, client.UserID)
			}
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

	hub.checkClients()
}

// util function.
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

// ListenAndBroadcast handles broadcasting
// pass in wg to wait for hub to shutdown before exiting application
func (hub *GameHub) ListenAndBroadcast(wg *sync.WaitGroup) {
	defer wg.Done()

	for !hub.IsShutdown {
		msg := <-hub.MessageBus

		if msg.Cmd == message.CmdShutdown {
			break
		}

		hub.Mutex.Lock()
		// process message
		var resp *message.GameMessage
		var isBroadcast = true
		if _, ok := hub.GameLoader.UserIDToPlayerMap[msg.Sender]; !ok { // non-player
			resp = message.GameErrorMessage(ErrMsgNotPlayer)
			isBroadcast = false
		} else if msg.Cmd == message.CmdChat {
			resp = msg
		} else {
			resp, isBroadcast = hub.GameLoader.HandleMessage(msg)
		}
		// broadcast
		if isBroadcast {
			for _, client := range hub.Clients {
				hub.sendMessageToClient(client, resp)
			}
		} else {
			hub.sendMessageToClient(hub.Clients[msg.Sender], resp)
		}
		hub.Mutex.Unlock()

		hub.checkClients()
	}
}

// if no clients left, game hub will shutdown
func (hub *GameHub) checkClients() {
	hub.Mutex.Lock()
	if !hub.IsShutdown {
		if len(hub.Clients) == 0 {
			// shutdown, but without lock
			hub.IsShutdown = true
			hub.GameLoader.SaveToDB()
			hub.MessageBus <- &message.GameMessage{Cmd: message.CmdShutdown} // trigger shutdown for the listening goroutine
			hub.OnShutdown()
		}
	}
	hub.Mutex.Unlock()
}

// ForceShutdown closes all client connection and stops the hub
func (hub *GameHub) ForceShutdown() {
	hub.Mutex.Lock()
	if !hub.IsShutdown {
		hub.IsShutdown = true
		hub.GameLoader.SaveToDB()
		hub.MessageBus <- &message.GameMessage{Cmd: message.CmdShutdown} // trigger shutdown for the listening goroutine
		for _, client := range hub.Clients {
			client.WS.Close()
			delete(hub.Clients, client.UserID)
		}
		// no need to call OnShutdown here
		// if called, it will deadlock anyway
	}
	hub.Mutex.Unlock()
}

// NewGameHub initializes a new game hub with the correct game id
func NewGameHub(gameID uint64, onShutdown func()) *GameHub {
	return &GameHub{
		GameID:     gameID,
		GameLoader: loader.NewGameLoader(gameID),
		Clients:    make(map[uint64]*GameClient),
		MessageBus: make(chan *message.GameMessage),
		Mutex:      sync.Mutex{},
		IsShutdown: false,
		OnShutdown: onShutdown,
	}
}
