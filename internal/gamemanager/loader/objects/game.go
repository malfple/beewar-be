package objects

import (
	"gitlab.com/otqee/otqee-be/internal/access/model"
	"gitlab.com/otqee/otqee-be/internal/gamemanager/loader/objects/units"
)

// Game is the main game object
// describes the map and everything else related to the game
type Game struct {
	ID          uint64
	Type        uint8
	Width       uint8
	Height      uint8
	PlayerCount uint8
	Terrain     [][]uint8
	Units       [][]*units.Unit
	MapID       uint64
	TurnCount   int32
	TurnPlayer  int8
}

// ToModel converts the current game object into a model.Game db model
func (game *Game) ToModel() *model.Game {
	// TODO: properly translate back to game model
	return nil
}

// NewGameFromModel creates a new game from model.Game db model
func NewGameFromModel(gameModel *model.Game) *Game {
	return &Game{
		ID:          gameModel.ID,
		Type:        gameModel.Type,
		Width:       gameModel.Width,
		Height:      gameModel.Height,
		PlayerCount: gameModel.PlayerCount,
		Terrain:     ModelToGameTerrain(gameModel.Width, gameModel.Height, gameModel.TerrainInfo),
		Units:       ModelToGameUnit(gameModel.Width, gameModel.Height, gameModel.UnitInfo),
		MapID:       gameModel.MapID,
		TurnCount:   gameModel.TurnCount,
		TurnPlayer:  gameModel.TurnPlayer,
	}
}

// TODO: move these functions to formatter package

// ModelToGameTerrain converts terrain info from model.Game to Game
func ModelToGameTerrain(width, height uint8, terrainInfo []byte) [][]uint8 {
	terrain := make([][]uint8, height)
	for i := uint8(0); i < height; i++ {
		terrain[i] = make([]uint8, width)
		for j := uint8(0); j < width; j++ {
			terrain[i][j] = terrainInfo[i*width+j]
		}
	}
	return terrain
}

// GameTerrainToModel converts terrain info from Game to model.Game
func GameTerrainToModel(width, height uint8, terrain [][]uint8) []byte {
	terrainInfo := make([]byte, int(width)*int(height))
	for i := uint8(0); i < height; i++ {
		for j := uint8(0); j < width; j++ {
			terrainInfo[i*width+j] = terrain[i][j]
		}
	}
	return terrainInfo
}

// ModelToGameUnit converts unit info from model.Game to Game
func ModelToGameUnit(width, height uint8, unitInfo []byte) [][]*units.Unit {
	_units := make([][]*units.Unit, height)
	for i := uint8(0); i < height; i++ {
		_units[i] = make([]*units.Unit, width)
		for j := uint8(0); j < width; j++ {
			_units[i][j] = nil
		}
	}
	// TODO: translate unit and assign to cells
	return _units
}

// GameUnitToModel converts unit info from Game to model.Game
func GameUnitToModel(width, height int8, _units [][]*units.Unit) []byte {
	return nil
}
