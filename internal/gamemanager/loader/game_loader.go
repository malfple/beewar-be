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
	Width        uint8
	Height       uint8
	PlayerCount  uint8
	Terrain      [][]uint8
	Units        [][]objects.Unit
	MapID        uint64
	TurnCount    int32
	TurnPlayer   int8
	TimeCreated  int64
	TimeModified int64
}

// NewGameLoader loads game by gameID and return the GameLoader object
func NewGameLoader(gameID uint64) *GameLoader {
	gameModel := access.QueryGameByID(gameID)
	if gameModel == nil {
		// the websocket handler should already handle this
		panic("loader: game is supposed to exist")
	}

	return &GameLoader{
		ID:           gameModel.ID,
		Type:         gameModel.Type,
		Width:        gameModel.Width,
		Height:       gameModel.Height,
		PlayerCount:  gameModel.PlayerCount,
		Terrain:      formatter.ModelToGameTerrain(gameModel.Width, gameModel.Height, gameModel.TerrainInfo),
		Units:        formatter.ModelToGameUnit(gameModel.Width, gameModel.Height, gameModel.UnitInfo),
		MapID:        gameModel.MapID,
		TurnCount:    gameModel.TurnCount,
		TurnPlayer:   gameModel.TurnPlayer,
		TimeCreated:  gameModel.TimeCreated,
		TimeModified: gameModel.TimeModified,
	}
}

// ToModel converts the current game object into a model.Game db model
func (gl *GameLoader) ToModel() *model.Game {
	return &model.Game{
		ID:           gl.ID,
		Type:         gl.Type,
		Width:        gl.Width,
		Height:       gl.Height,
		PlayerCount:  gl.PlayerCount,
		TerrainInfo:  formatter.GameTerrainToModel(gl.Width, gl.Height, gl.Terrain),
		UnitInfo:     formatter.GameUnitToModel(gl.Width, gl.Height, gl.Units),
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
