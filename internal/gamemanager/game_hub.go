package gamemanager

import (
	"github.com/gorilla/websocket"
	"sync"
)

// the hub is responsible for broadcasting messages to the clients

// GameHub is the central broadcasting service for a game
type GameHub struct {
	GameID     int64
	Clients    map[*GameClient]bool
	MessageBus chan string
	Mutex      sync.Mutex // this lock is used to make sure broadcast and client register/unreg is synced
	IsShutdown bool
	OnShutdown func() // called when not forced to shutdown
}

// RegisterClient registers the client to this hub
func (hub *GameHub) RegisterClient(client *GameClient) {
	hub.Mutex.Lock()
	if !hub.IsShutdown {
		hub.Clients[client] = true
	}
	hub.Mutex.Unlock()
}

// UnregisterClient unregisters the client
func (hub *GameHub) UnregisterClient(client *GameClient) {
	hub.Mutex.Lock()
	if !hub.IsShutdown {
		delete(hub.Clients, client)
	}
	hub.Mutex.Unlock()

	hub.checkClients()
}

// ListenAndBroadcast handles broadcasting
// pass in wg to wait for hub to shutdown before exiting application
func (hub *GameHub) ListenAndBroadcast(wg *sync.WaitGroup) {
	defer wg.Done()

	for !hub.IsShutdown {
		hub.checkClients()

		msg := <-hub.MessageBus

		if msg == "shutdown" {
			break
		}

		hub.Mutex.Lock()
		for client := range hub.Clients {
			err := client.WS.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				client.WS.Close()
				delete(hub.Clients, client)
			}
		}
		hub.Mutex.Unlock()
	}
}

// if no clients left, game hub will shutdown
func (hub *GameHub) checkClients() {
	hub.Mutex.Lock()
	if !hub.IsShutdown {
		if len(hub.Clients) == 0 {
			// shutdown, but without lock
			hub.IsShutdown = true
			hub.MessageBus <- "shutdown" // trigger shutdown for the listening goroutine
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
		hub.MessageBus <- "shutdown" // trigger shutdown for the listening goroutine
		for client := range hub.Clients {
			client.WS.Close()
			delete(hub.Clients, client)
		}
		// no need to call OnShutdown here
		// if called, it will deadlock anyway
	}
	hub.Mutex.Unlock()
}

// NewGameHub initializes a new game hub with the correct game id
func NewGameHub(gameID int64, onShutdown func()) *GameHub {
	return &GameHub{
		GameID:     gameID,
		Clients:    make(map[*GameClient]bool),
		MessageBus: make(chan string),
		Mutex:      sync.Mutex{},
		IsShutdown: false,
		OnShutdown: onShutdown,
	}
}
