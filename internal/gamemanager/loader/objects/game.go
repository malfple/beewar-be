package objects

import "gitlab.com/otqee/otqee-be/internal/access/model"

// Game is the main game object
type Game struct {
	ID          int64
	Type        int8
	Width       int8
	Height      int8
	PlayerCount int8
	MapID       int64
	TurnCount   int32
	TurnPlayer  int8
}

// NewGameFromDBModel creates a new game from model.Game db model
func NewGameFromDBModel(gameModel *model.Game) *Game {
	// TODO: add terrain and unit
	return &Game{
		ID:          gameModel.ID,
		Type:        gameModel.Type,
		Width:       gameModel.Width,
		Height:      gameModel.Height,
		PlayerCount: gameModel.PlayerCount,
		MapID:       gameModel.MapID,
		TurnCount:   gameModel.TurnCount,
		TurnPlayer:  gameModel.TurnPlayer,
	}
}

// ToDBModel converts the current game object into a model.Game db model
func (game *Game) ToDBModel() *model.Game {
	// TODO: properly translate back to game model
	return nil
}
