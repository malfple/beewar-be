package gamemanager

import (
	"gitlab.com/beewar/beewar-be/internal/logger"
	"go.uber.org/zap"
	"sync"
)

var gameHubStore = make(map[uint64]*GameHub)
var gameHubStoreLock sync.Mutex
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
		hub.Shutdown()
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

// StartClientSession will register a given client to the game hub with the given gameID.
// it will initialize the hub if it is not yet initialized.
// returns error if initialization or registration fails.
func StartClientSession(client GameClient, gameID uint64) error {
	gameHubStoreLock.Lock()
	defer gameHubStoreLock.Unlock()

	// game hub fetch / creation
	hub, ok := gameHubStore[gameID]
	if !ok {
		var err error
		logger.GetLogger().Debug("game manager: open new game hub", zap.Uint64("game_id", gameID))
		hub, err = NewGameHub(gameID)
		if err != nil {
			return err
		}

		gameHubWG.Add(1)
		// goroutine for the hub to do its job
		go hub.ListenAndBroadcast(&gameHubWG)

		gameHubStore[gameID] = hub
	}

	// error is usually caused by duplicate client, so it's okay to just exit prematurely
	err := hub.RegisterClient(client)
	if err != nil {
		return err
	}

	client.SetHub(hub)
	return nil
}

// EndClientSession will unregister a given client from its game hub.
// It also checks game hub conditions to allow hub shutdown.
func EndClientSession(client GameClient) {
	gameHubStoreLock.Lock()
	defer gameHubStoreLock.Unlock()

	hub := client.GetHub()
	hub.UnregisterClient(client)

	// shutdown hub and remove from store if no clients left
	if len(hub.Clients) == 0 {
		logger.GetLogger().Debug("game manager: close game hub", zap.Uint64("game_id", hub.GameID))
		delete(gameHubStore, hub.GameID)
		hub.Shutdown()
	}
}
