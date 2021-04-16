package loader

import (
	"errors"
	"gitlab.com/beewar/beewar-be/internal/access"
	"gitlab.com/beewar/beewar-be/internal/access/model"
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager/formatter"
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager/loader/gridengine"
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager/message"
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager/objects"
	"gitlab.com/beewar/beewar-be/internal/logger"
	"go.uber.org/zap"
)

const (
	// ErrMsgGameNotStarted is returned when a message is received and the game is not yet started
	ErrMsgGameNotStarted = "game not yet started"
	// ErrMsgGameEnded is returned when a message is received and the game is already ended
	ErrMsgGameEnded = "game already ended"
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
	// ErrMsgInvalidAttack is returned when an attack is invalid
	ErrMsgInvalidAttack = "invalid attack. maybe your target is out of range"
)

var (
	errGameDoesNotExist = errors.New("game does not exist")
	errGameInPicking    = errors.New("game is still in picking phase")
)

const (
	// GameStatusPicking is a game state
	GameStatusPicking = 0
	// GameStatusOngoing is a game state
	GameStatusOngoing = 1
	// GameStatusEnded is a game state
	GameStatusEnded = 2
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
	Password          string
	Status            int8  // indicates game status
	TurnCount         int32 // turns start from 1, defined in migration
	TurnPlayer        int   // players are numbered 1..PlayerCount.
	TimeCreated       int64
	TimeModified      int64
	GridEngine        *gridengine.GridEngine
	GameUsers         []*model.GameUser
	Users             []*model.User
	UserIDToPlayerMap map[uint64]int
}

// NewGameLoader loads game by gameID and return the GameLoader object
func NewGameLoader(gameID uint64) (*GameLoader, error) {
	gameModel := access.QueryGameByID(gameID)
	if gameModel == nil {
		return nil, errGameDoesNotExist
	}
	gameUsers := access.QueryGameUsersByGameID(gameID)
	if gameModel.Status == GameStatusPicking {
		if int(gameModel.PlayerCount) == len(gameUsers) { // auto-start
			gameModel.Status = GameStatusOngoing
		} else {
			return nil, errGameInPicking
		}
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
		Password:     gameModel.Password,
		Status:       gameModel.Status,
		TurnCount:    gameModel.TurnCount,
		TurnPlayer:   int(gameModel.TurnPlayer),
		TimeCreated:  gameModel.TimeCreated,
		TimeModified: gameModel.TimeModified,
	}
	// create grid engine
	gameLoader.GridEngine = gridengine.NewGridEngine(
		gameLoader.Width,
		gameLoader.Height,
		&gameLoader.Terrain,
		&gameLoader.Units)
	// load players
	gameLoader.GameUsers = gameUsers
	userIDs := make([]uint64, len(gameLoader.GameUsers))
	for i, gu := range gameLoader.GameUsers {
		userIDs[i] = gu.UserID
	}
	gameLoader.Users = access.QueryUsersByID(userIDs)
	for i := range gameLoader.Users {
		gameLoader.Users[i].Password = "nope."
	}
	// make reverse map
	gameLoader.UserIDToPlayerMap = make(map[uint64]int)
	for _, user := range gameLoader.GameUsers {
		gameLoader.UserIDToPlayerMap[user.UserID] = int(user.PlayerOrder)
	}

	return gameLoader, nil
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
		Password:     gl.Password,
		Status:       gl.Status,
		TurnCount:    gl.TurnCount,
		TurnPlayer:   int8(gl.TurnPlayer),
		TimeCreated:  gl.TimeCreated,
		TimeModified: gl.TimeModified,
	}
}

// SaveToDB saves the current game object to db
func (gl *GameLoader) SaveToDB() error {
	gameModel := gl.ToModel()
	if err := formatter.ValidateUnitInfo(gameModel.Height, gameModel.Width, gameModel.UnitInfo); err != nil {
		logger.GetLogger().Error("loader: fail unit info validation when saving", zap.Error(err))
		return err
	}
	if err := access.UpdateGameAndGameUser(gameModel, gl.GameUsers); err != nil {
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
			User:        gl.Users[i],
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

func (gl *GameLoader) checkGameEnd() {
	playersLeft := 0
	for _, gu := range gl.GameUsers {
		if gu.FinalTurns == 0 {
			playersLeft++
		}
	}
	if playersLeft <= 1 {
		// game ended
		for i := range gl.GameUsers {
			gl.assignPlayerRank(i + 1)
		}
		gl.Status = GameStatusEnded
	}
}

// end current player turn and start next player turn
func (gl *GameLoader) nextTurn() {
	prevPlayer := gl.TurnPlayer
	for {
		gl.TurnPlayer++
		if gl.TurnPlayer > gl.PlayerCount {
			gl.TurnCount++
			gl.TurnPlayer = 1
		}
		if gl.GameUsers[gl.TurnPlayer-1].FinalTurns == 0 { // player not defeated
			break
		}
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
	// just in case, also check if game ends
	gl.checkGameEnd()
}

// check if a unit is still alive. Also checks player defeat condition. If you dies, you are defeated.
func (gl *GameLoader) checkUnitAlive(y, x int) {
	if gl.Units[y][x].GetUnitHP() == 0 {
		if gl.Units[y][x].GetUnitType() == objects.UnitTypeYou {
			// player is defeated -> assign rank and turns lasted
			gl.assignPlayerRank(gl.Units[y][x].GetUnitOwner())
			// immediately check if game ends
			gl.checkGameEnd()
		}
		gl.Units[y][x] = nil
	}
}

// done when player is defeated or when game is finished
func (gl *GameLoader) assignPlayerRank(player int) {
	if gl.GameUsers[player-1].FinalTurns != 0 {
		return
	}
	for _, gu := range gl.GameUsers {
		if gu.FinalTurns == 0 { // player not yet defeated
			gl.GameUsers[player-1].FinalRank++
		}
	}
	gl.GameUsers[player-1].FinalTurns = gl.TurnCount
}

// HandleMessage handles game related message
// returns the message and a boolean value whether the message should be broadcasted (true = broadcast)
func (gl *GameLoader) HandleMessage(msg *message.GameMessage) (*message.GameMessage, bool) {
	if gl.Status == GameStatusPicking {
		return message.GameErrorMessage(ErrMsgGameNotStarted), false
	}
	if gl.Status == GameStatusEnded {
		return message.GameErrorMessage(ErrMsgGameEnded), false
	}

	// only current player can do stuff
	if msg.Sender != gl.GameUsers[gl.TurnPlayer-1].UserID {
		return message.GameErrorMessage(ErrMsgNotYetTurn), false
	}

	switch msg.Cmd {
	case message.CmdUnitMove:
		return gl.handleUnitMove(msg)
	case message.CmdUnitAttack:
		return gl.handleUnitAttack(msg)
	case message.CmdUnitMoveAndAttack:
		return gl.handleUnitMoveAndAttack(msg)
	case message.CmdEndTurn:
		gl.nextTurn()
		return msg, true
	}
	panic("panic game loader handle message: cmd not allowed")
}
