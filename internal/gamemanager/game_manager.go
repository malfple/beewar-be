package gamemanager

import (
	"github.com/gorilla/websocket"
	"gitlab.com/otqee/otqee-be/internal/logger"
	"go.uber.org/zap"
	"sync"
)

var gameHubStore = make(map[int64]*GameHub)
var gameHubStoreLock sync.RWMutex
var wg sync.WaitGroup

// InitGameManager starts the game manager
func InitGameManager() {
	logger.GetLogger().Info("game manager: init")
	// actually do nothing
}

// ShutdownGameManager stops the game manager and closes all connections
func ShutdownGameManager() {
	logger.GetLogger().Info("game manager: shutting down hubs", zap.Int("hub_count", len(gameHubStore)))
	gameHubStoreLock.Lock()
	for gameID, hub := range gameHubStore {
		hub.ForceShutdown()
		delete(gameHubStore, gameID)
	}
	wg.Wait()
	gameHubStoreLock.Unlock()
	logger.GetLogger().Info("game manager: shutdown")
}

// GetGameHub returns the game hub with the corresponding game id
// it will initialize the hub if it is not yet initialized
func GetGameHub(gameID int64) *GameHub {
	gameHubStoreLock.RLock()
	hub, ok := gameHubStore[gameID]
	gameHubStoreLock.RUnlock()
	if ok {
		return hub
	}

	logger.GetLogger().Debug("game manager: open new game hub", zap.Int64("game_id", gameID))
	hub = NewGameHub(gameID, func() {
		logger.GetLogger().Debug("game manager: close game hub", zap.Int64("game_id", gameID))
		gameHubStoreLock.Lock()
		delete(gameHubStore, gameID)
		gameHubStoreLock.Unlock()
	})
	wg.Add(1)
	// goroutine for the hub to do its job
	go hub.ListenAndBroadcast(&wg)

	gameHubStoreLock.Lock()
	gameHubStore[gameID] = hub
	gameHubStoreLock.Unlock()
	return hub
}

// NewGameClientByID creates a new client and connects to the hub by game id
func NewGameClientByID(ws *websocket.Conn, gameID int64) *GameClient {
	return NewGameClient(ws, GetGameHub(gameID))
}
