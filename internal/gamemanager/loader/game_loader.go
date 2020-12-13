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
	// ErrMsgNotYetTurn is returned when it is not yet the turn for the player who sends the message
	ErrMsgNotYetTurn = "not your turn yet"
	// ErrMsgInvalidPos is returned if given position does not have a unit or is out of bound
	ErrMsgInvalidPos = "you tried to move a non-existent unit or position is out of map"
	// ErrMsgUnitNotOwned is returned when the player tries to move a unit they do not own
	ErrMsgUnitNotOwned = "you tried to move a unit you do not own"
	// ErrMsgUnitCmdNotAllowed is returned if a cmd cannot be used on the unit. See CmdWhitelist
	ErrMsgUnitCmdNotAllowed = "the cmd cannot be used on that unit"
	// ErrMsgUnitAlreadyMoved is returned when the unit already moved or attacked
	ErrMsgUnitAlreadyMoved = "unit already moved or attacked"
	// ErrMsgInvalidMove is returned when a move is invalid
	ErrMsgInvalidMove = "invalid move duh"
)

// GameLoader loads game from db and perform game tasks.
// also supports saving back to db on demand (only the game -> game_tab).
// game loader is not concurrent safe, and the caller needs to handle this with locks
// players information (game_user_tab) will be saved to db when a player is defeated
type GameLoader struct {
	ID                uint64 // this is game id
	Type              uint8
	Height            int
	Width             int
	PlayerCount       int
	Terrain           [][]int
	Units             [][]objects.Unit
	MapID             uint64
	TurnCount         int32 // turns start from 1, defined in migration
	TurnPlayer        int   // players are numbered 1..PlayerCount
	TimeCreated       int64
	TimeModified      int64
	GridEngine        *GridEngine
	GameUsers         []*model.GameUser
	UserIDToPlayerMap map[uint64]int
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
	// make reverse map
	gameLoader.UserIDToPlayerMap = make(map[uint64]int)
	for _, user := range gameLoader.GameUsers {
		gameLoader.UserIDToPlayerMap[user.UserID] = int(user.PlayerOrder)
	}

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

// end current player turn and start next player turn
func (gl *GameLoader) nextTurn() {
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

// validate functions return empty string if no errors

// validates of unit is owned by the userID given.
// also validates position inside map, and if a unit exists in the given position
func (gl *GameLoader) validateUnitOwned(userID uint64, y, x int) string {
	if y < 0 || y > gl.Height || x < 0 || x > gl.Width {
		return ErrMsgInvalidPos
	}
	if gl.Units[y][x] == nil {
		return ErrMsgInvalidPos
	}
	// player doesn't own the unit
	if gl.UserIDToPlayerMap[userID] != gl.Units[y][x].GetUnitOwner() {
		return ErrMsgUnitNotOwned
	}

	return ""
}

// HandleMessage handles game related message
// returns the message and a boolean value whether the message should be broadcasted (true = broadcast)
func (gl *GameLoader) HandleMessage(msg *message.GameMessage) (*message.GameMessage, bool) {
	// only current player can do stuff
	if msg.Sender != gl.GameUsers[gl.TurnPlayer-1].UserID {
		return message.GameErrorMessage(ErrMsgNotYetTurn), false
	}

	switch msg.Cmd {
	case message.CmdUnitMove:
		data := msg.Data.(*message.UnitMoveMessageData)
		if errMsg := gl.validateUnitOwned(msg.Sender, data.Y1, data.X1); len(errMsg) > 0 {
			return message.GameErrorMessage(errMsg), false
		}
		if _, ok := CmdWhitelistUnitMove[gl.Units[data.Y1][data.X1].GetUnitType()]; !ok {
			return message.GameErrorMessage(ErrMsgUnitCmdNotAllowed), false
		}
		if gl.Units[data.Y1][data.X1].GetUnitStateBit(objects.UnitStateBitMoved) {
			return message.GameErrorMessage(ErrMsgUnitAlreadyMoved), false
		}
		if !gl.GridEngine.ValidateMove(data.Y1, data.X1, data.Y2, data.X2) {
			return message.GameErrorMessage(ErrMsgInvalidMove), false
		}
		gl.Units[data.Y2][data.X2] = gl.Units[data.Y1][data.X1]
		gl.Units[data.Y2][data.X2].ToggleUnitStateBit(objects.UnitStateBitMoved)
		gl.Units[data.Y1][data.X1] = nil
		// TODO: make EVENT message and return that instead
		return msg, true
	case message.CmdEndTurn:
		gl.nextTurn()
		return msg, true
	}
	panic("panic game loader handle message: cmd not allowed")
}
