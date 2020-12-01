package loader

import (
	"gitlab.com/otqee/otqee-be/internal/access"
	"gitlab.com/otqee/otqee-be/internal/gamemanager/loader/objects"
	"gitlab.com/otqee/otqee-be/internal/logger"
	"go.uber.org/zap"
)

// GameLoader loads game from db and perform game tasks.
// also supports saving back to db on demand.
// game loader is not concurrent safe, and the caller needs to handle this with locks
type GameLoader struct {
	Game *objects.Game
}

// NewGameLoader loads game by gameID and return the GameLoader object
func NewGameLoader(gameID uint64) *GameLoader {
	game := access.QueryGameByID(gameID)
	if game == nil {
		// the websocket handler should already handle this
		panic("loader: game is supposed to exist")
	}

	return &GameLoader{
		Game: objects.NewGameFromModel(game),
	}
}

// SaveToDB saves the current game object to db
func (gl *GameLoader) SaveToDB() error {
	gameModel := gl.Game.ToModel()
	if err := access.UpdateGame(gameModel); err != nil {
		logger.GetLogger().Error("loader: error save game to db", zap.Error(err))
		return err
	}
	return nil
}
