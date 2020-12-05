package loader

import (
	"gitlab.com/otqee/otqee-be/internal/access"
	"gitlab.com/otqee/otqee-be/internal/access/formatter"
	"gitlab.com/otqee/otqee-be/internal/access/formatter/objects"
	"gitlab.com/otqee/otqee-be/internal/access/model"
	"gitlab.com/otqee/otqee-be/internal/logger"
	"go.uber.org/zap"
)

// GameLoader loads game from db and perform game tasks.
// also supports saving back to db on demand.
// game loader is not concurrent safe, and the caller needs to handle this with locks
type GameLoader struct {
	ID           uint64
	Type         uint8
	Height       int
	Width        int
	PlayerCount  uint8
	Terrain      [][]int
	Units        [][]objects.Unit
	MapID        uint64
	TurnCount    int32
	TurnPlayer   int8
	TimeCreated  int64
	TimeModified int64
	GridEngine   *GridEngine
}

// NewGameLoader loads game by gameID and return the GameLoader object
func NewGameLoader(gameID uint64) *GameLoader {
	gameModel := access.QueryGameByID(gameID)
	if gameModel == nil {
		// the websocket handler should already handle this
		panic("loader: game is supposed to exist")
	}

	gameLoader := &GameLoader{
		ID:           gameModel.ID,
		Type:         gameModel.Type,
		Height:       gameModel.Height,
		Width:        gameModel.Width,
		PlayerCount:  gameModel.PlayerCount,
		Terrain:      formatter.ModelToGameTerrain(gameModel.Height, gameModel.Width, gameModel.TerrainInfo),
		Units:        formatter.ModelToGameUnit(gameModel.Height, gameModel.Width, gameModel.UnitInfo),
		MapID:        gameModel.MapID,
		TurnCount:    gameModel.TurnCount,
		TurnPlayer:   gameModel.TurnPlayer,
		TimeCreated:  gameModel.TimeCreated,
		TimeModified: gameModel.TimeModified,
	}

	gameLoader.GridEngine = NewGridEngine(
		gameLoader.Width,
		gameLoader.Height,
		&gameLoader.Terrain,
		&gameLoader.Units)

	return gameLoader
}

// ToModel converts the current game object into a model.Game db model
func (gl *GameLoader) ToModel() *model.Game {
	return &model.Game{
		ID:           gl.ID,
		Type:         gl.Type,
		Height:       gl.Height,
		Width:        gl.Width,
		PlayerCount:  gl.PlayerCount,
		TerrainInfo:  formatter.GameTerrainToModel(gl.Height, gl.Width, gl.Terrain),
		UnitInfo:     formatter.GameUnitToModel(gl.Height, gl.Width, gl.Units),
		MapID:        gl.MapID,
		TurnCount:    gl.TurnCount,
		TurnPlayer:   gl.TurnPlayer,
		TimeCreated:  gl.TimeCreated,
		TimeModified: gl.TimeModified,
	}
}

// SaveToDB saves the current game object to db
func (gl *GameLoader) SaveToDB() error {
	gameModel := gl.ToModel()
	if err := access.UpdateGame(gameModel); err != nil {
		logger.GetLogger().Error("loader: error save game to db", zap.Error(err))
		return err
	}
	return nil
}
