package gamemanager

import (
	"gitlab.com/beewar/beewar-be/internal/logger"
	"go.uber.org/zap"
	"sync"
)

var gameHubStore = make(map[uint64]*GameHub)
var gameHubStoreLock sync.RWMutex
var gameHubWG sync.WaitGroup

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
		logger.GetLogger().Debug("game manager: force shutdown hub", zap.Uint64("game_id", hub.GameID))
		hub.ForceShutdown()
		delete(gameHubStore, gameID)
	}
	gameHubWG.Wait()
	gameHubStoreLock.Unlock()
	logger.GetLogger().Info("game manager: shutdown")
}

// GetHubCount returns the number of active hubs
func GetHubCount() int {
	return len(gameHubStore)
}

// GetGameHub returns the game hub with the corresponding game id.
// it will initialize the hub if it is not yet initialized.
// returns error if initialization fails.
func GetGameHub(gameID uint64) (*GameHub, error) {
	gameHubStoreLock.RLock()
	hub, ok := gameHubStore[gameID]
	gameHubStoreLock.RUnlock()
	if ok {
		return hub, nil
	}

	var err error
	logger.GetLogger().Debug("game manager: open new game hub", zap.Uint64("game_id", gameID))
	hub, err = NewGameHub(gameID, func() {
		logger.GetLogger().Debug("game manager: close game hub", zap.Uint64("game_id", gameID))
		gameHubStoreLock.Lock()
		delete(gameHubStore, gameID)
		gameHubStoreLock.Unlock()
	})
	if err != nil {
		return nil, err
	}
	gameHubWG.Add(1)
	// goroutine for the hub to do its job
	go hub.ListenAndBroadcast(&gameHubWG)

	gameHubStoreLock.Lock()
	gameHubStore[gameID] = hub
	gameHubStoreLock.Unlock()
	return hub, nil
}
