package loader

import (
	"gitlab.com/otqee/otqee-be/internal/access"
	"gitlab.com/otqee/otqee-be/internal/access/formatter"
	"gitlab.com/otqee/otqee-be/internal/access/formatter/objects"
	"gitlab.com/otqee/otqee-be/internal/access/model"
	"gitlab.com/otqee/otqee-be/internal/gamemanager/message"
	"gitlab.com/otqee/otqee-be/internal/logger"
	"go.uber.org/zap"
)

const (
	// ErrMsgUnauthorized is returned when the player who sends the message is not authorized
	ErrMsgUnauthorized = "player is not authorized to perform the cmd"
)

// GameLoader loads game from db and perform game tasks.
// also supports saving back to db on demand (only the game -> game_tab).
// game loader is not concurrent safe, and the caller needs to handle this with locks
// players information (game_user_tab) will be saved to db when a player is defeated
type GameLoader struct {
	ID           uint64 // this is game id
	Type         uint8
	Height       int
	Width        int
	PlayerCount  int
	Terrain      [][]int
	Units        [][]objects.Unit
	MapID        uint64
	TurnCount    int32 // turns start from 1, defined in migration
	TurnPlayer   int   // players are numbered 1..PlayerCount
	TimeCreated  int64
	TimeModified int64
	GridEngine   *GridEngine
	GameUsers    []*model.GameUser
}

// NewGameLoader loads game by gameID and return the GameLoader object
func NewGameLoader(gameID uint64) *GameLoader {
	gameModel := access.QueryGameByID(gameID)
	if gameModel == nil {
		// the websocket handler should already handle this
		panic("loader: game is supposed to exist")
	}
	// load main fields from game model
	gameLoader := &GameLoader{
		ID:           gameModel.ID,
		Type:         gameModel.Type,
		Height:       gameModel.Height,
		Width:        gameModel.Width,
		PlayerCount:  int(gameModel.PlayerCount),
		Terrain:      formatter.ModelToGameTerrain(gameModel.Height, gameModel.Width, gameModel.TerrainInfo),
		Units:        formatter.ModelToGameUnit(gameModel.Height, gameModel.Width, gameModel.UnitInfo),
		MapID:        gameModel.MapID,
		TurnCount:    gameModel.TurnCount,
		TurnPlayer:   int(gameModel.TurnPlayer),
		TimeCreated:  gameModel.TimeCreated,
		TimeModified: gameModel.TimeModified,
	}
	// create grid engine
	gameLoader.GridEngine = NewGridEngine(
		gameLoader.Width,
		gameLoader.Height,
		&gameLoader.Terrain,
		&gameLoader.Units)
	// load players
	gameLoader.GameUsers = access.QueryUsersLinkedToGame(gameLoader.ID)

	return gameLoader
}

// ToModel converts the current game object into a model.Game db model
func (gl *GameLoader) ToModel() *model.Game {
	return &model.Game{
		ID:           gl.ID,
		Type:         gl.Type,
		Height:       gl.Height,
		Width:        gl.Width,
		PlayerCount:  uint8(gl.PlayerCount),
		TerrainInfo:  formatter.GameTerrainToModel(gl.Height, gl.Width, gl.Terrain),
		UnitInfo:     formatter.GameUnitToModel(gl.Height, gl.Width, gl.Units),
		MapID:        gl.MapID,
		TurnCount:    gl.TurnCount,
		TurnPlayer:   int8(gl.TurnPlayer),
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

// GameData returns the game data message of the loaded game
func (gl *GameLoader) GameData() *message.GameMessage {
	players := make([]*message.Player, len(gl.GameUsers))
	for i, gameUser := range gl.GameUsers {
		players[i] = &message.Player{
			UserID:      gameUser.UserID,
			PlayerOrder: gameUser.PlayerOrder,
			FinalRank:   gameUser.FinalRank,
			FinalTurns:  gameUser.FinalTurns,
		}
	}
	return &message.GameMessage{
		Cmd: message.CmdGameData,
		Data: &message.GameDataMessageData{
			Game:    gl.ToModel(),
			Players: players,
		},
	}
}

// end current player turn
func (gl *GameLoader) endTurn() {
	prevPlayer := gl.TurnPlayer
	gl.TurnPlayer++
	if gl.TurnPlayer > gl.PlayerCount {
		gl.TurnCount++
		gl.TurnPlayer = 1
	}
	// unit states
	for i := 0; i < gl.Height; i++ {
		for j := 0; j < gl.Width; j++ {
			if gl.Units[i][j] == nil {
				continue
			}
			if gl.Units[i][j].GetUnitOwner() == prevPlayer {
				gl.Units[i][j].EndTurn()
			} else if gl.Units[i][j].GetUnitOwner() == gl.TurnPlayer {
				gl.Units[i][j].StartTurn()
			}
		}
	}
}

// HandleMessage handles game related message
// returns the message and a boolean value whether the message should be broadcasted (true = broadcast)
func (gl *GameLoader) HandleMessage(msg *message.GameMessage) (*message.GameMessage, bool) {
	switch msg.Cmd {
	case message.CmdUnitMove:
		// TODO: implement
	case message.CmdEndTurn:
		// only current user can end current turn
		if msg.Sender != gl.GameUsers[gl.TurnPlayer-1].UserID {
			return message.GameErrorMessage(ErrMsgUnauthorized), false
		}
		gl.endTurn()
		return msg, true
	}
	panic("panic game loader handle message: cmd not allowed")
}
