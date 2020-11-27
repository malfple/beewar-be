package loader

import (
	"gitlab.com/otqee/otqee-be/internal/access"
	"gitlab.com/otqee/otqee-be/internal/access/model"
	"gitlab.com/otqee/otqee-be/internal/gamemanager/loader/objects"
)

// GameLoader loads game from db and perform game tasks.
// also supports saving back to db on demand.
// by default, GameLoader will save back to db on end of turn.
// game loader is not concurrent safe, and the caller needs to handle this with locks
type GameLoader struct {
	Game *objects.Game
}

// NewGameLoader loads game by gameID and return the GameLoader object and model.Game db model
func NewGameLoader(gameID int64) (*GameLoader, *model.Game) {
	game := access.QueryGameByID(gameID)
	if game == nil {
		// the websocket handler should already handle this
		panic("game is supposed to exist")
	}

	return &GameLoader{
		Game: objects.NewGameFromDBModel(game),
	}, game
}

// TODO: add save function (need access layer first
