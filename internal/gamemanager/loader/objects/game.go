package objects

import (
	"gitlab.com/otqee/otqee-be/internal/access/formatter"
	"gitlab.com/otqee/otqee-be/internal/access/model"
	"gitlab.com/otqee/otqee-be/internal/gamemanager/loader/objects/units"
)

// Game is the main game object
// describes the map and everything else related to the game
type Game struct {
	ID           uint64
	Type         uint8
	Width        uint8
	Height       uint8
	PlayerCount  uint8
	Terrain      [][]uint8
	Units        [][]units.Unit
	MapID        uint64
	TurnCount    int32
	TurnPlayer   int8
	TimeCreated  int64
	TimeModified int64
}

// NewGameFromModel creates a new game from model.Game db model
func NewGameFromModel(gameModel *model.Game) *Game {
	return &Game{
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
func (game *Game) ToModel() *model.Game {
	return &model.Game{
		ID:           game.ID,
		Type:         game.Type,
		Width:        game.Width,
		Height:       game.Height,
		PlayerCount:  game.PlayerCount,
		TerrainInfo:  formatter.GameTerrainToModel(game.Width, game.Height, game.Terrain),
		UnitInfo:     formatter.GameUnitToModel(game.Width, game.Height, game.Units),
		MapID:        game.MapID,
		TurnCount:    game.TurnCount,
		TurnPlayer:   game.TurnPlayer,
		TimeCreated:  game.TimeCreated,
		TimeModified: game.TimeModified,
	}
}
