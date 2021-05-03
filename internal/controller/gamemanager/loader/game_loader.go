package loader

import (
	"errors"
	"gitlab.com/beewar/beewar-be/internal/access"
	"gitlab.com/beewar/beewar-be/internal/access/model"
	"gitlab.com/beewar/beewar-be/internal/controller/formatter"
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager/loader/gridengine"
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager/message"
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager/objects"
	"gitlab.com/beewar/beewar-be/internal/logger"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

const (
	// joining error messages
	errMsgGameNotInPicking   = "game picking phase is over"
	errMsgPlayerOrderInvalid = "invalid slot/player_order given"
	errMsgPlayerOrderTaken   = "that slot/player_order is already taken"
	errMsgAlreadyJoined      = "you have already joined this game"
	errMsgWrongPassword      = "wrong password"
	// general validation errors
	errMsgGameNotStarted = "game not yet started"
	errMsgGameEnded      = "game already ended"
	errMsgNotYetTurn     = "not your turn yet"
	// unit cmd validation errors
	errMsgInvalidPos        = "you tried to move a non-existent unit or position is out of map"
	errMsgUnitNotOwned      = "you tried to move a unit you do not own"
	errMsgUnitCmdNotAllowed = "the cmd cannot be used on that unit" // check cmdwhitelist
	errMsgUnitAlreadyMoved  = "unit already moved or attacked"
	errMsgInvalidMove       = "invalid move duh"
	errMsgInvalidAttack     = "invalid attack. maybe your target is out of range"
)

var (
	errGameDoesNotExist = errors.New("game does not exist")
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
	gameUsers := access.QueryGameUsersByGameID(gameID)
	gameLoader.GameUsers = padGameUsers(gameUsers, gameLoader.PlayerCount)
	userIDs := make([]uint64, len(gameLoader.GameUsers))
	for i, gu := range gameLoader.GameUsers {
		userIDs[i] = gu.UserID
	}
	gameLoader.Users = access.QueryUsersByID(userIDs)
	for i := range gameLoader.Users {
		if gameLoader.Users[i] != nil {
			gameLoader.Users[i].Password = "nope."
		}
	}
	// make reverse map
	gameLoader.UserIDToPlayerMap = make(map[uint64]int)
	for _, user := range gameLoader.GameUsers {
		gameLoader.UserIDToPlayerMap[user.UserID] = int(user.PlayerOrder)
	}

	gameLoader.checkGameStart()

	return gameLoader, nil
}

// ToModel converts the current game object into a model.Game db model
func (gl *GameLoader) ToModel() *model.Game {
	game := &model.Game{
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
	if game.Password != "" { // only mask if password exists
		game.Password = "masked haha" // doesn't matter even on db save because cannot update password anyway
	}
	return game
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

// handles player join
func (gl *GameLoader) handleJoin(msg *message.GameMessage) (*message.GameMessage, bool) {
	if gl.Status != GameStatusPicking {
		return message.GameErrorMessage(errMsgGameNotInPicking), false
	}
	data := msg.Data.(*message.JoinMessageData)
	if data.PlayerOrder < 1 || int(data.PlayerOrder) > gl.PlayerCount {
		return message.GameErrorMessage(errMsgPlayerOrderInvalid), false
	}
	if gl.GameUsers[data.PlayerOrder-1].UserID != 0 { // taken haha
		return message.GameErrorMessage(errMsgPlayerOrderTaken), false
	}
	if _, ok := gl.UserIDToPlayerMap[msg.Sender]; ok {
		return message.GameErrorMessage(errMsgAlreadyJoined), false
	}
	if gl.Password != "" {
		if err := bcrypt.CompareHashAndPassword([]byte(gl.Password), []byte(data.Password)); err != nil {
			return message.GameErrorMessage(errMsgWrongPassword), false
		}
	}
	// pass join validation
	if err := access.CreateGameUser(gl.ID, msg.Sender, data.PlayerOrder); err != nil {
		return message.GameErrorMessage(err.Error()), false
	}
	gl.GameUsers[data.PlayerOrder-1] = access.QueryGameUser(gl.ID, msg.Sender)
	gl.Users[data.PlayerOrder-1] = access.QueryUsersByID([]uint64{msg.Sender})[0]
	gl.Users[data.PlayerOrder-1].Password = "nope..."
	gl.UserIDToPlayerMap[msg.Sender] = int(data.PlayerOrder)
	gl.checkGameStart()
	return &message.GameMessage{
		Cmd:    msg.Cmd,
		Sender: msg.Sender,
		Data: &message.JoinMessageDataExt{
			Player: &message.Player{
				UserID:      msg.Sender,
				PlayerOrder: data.PlayerOrder,
				FinalRank:   0,
				FinalTurns:  0,
				User:        gl.Users[data.PlayerOrder-1],
			},
		},
	}, true
}

func (gl *GameLoader) checkGameStart() {
	if gl.Status != GameStatusPicking {
		return
	}
	for _, gu := range gl.GameUsers {
		if gu.UserID == 0 {
			// empty slot
			return
		}
	}
	logger.GetLogger().Debug("loader: game started", zap.Uint64("game_id", gl.ID))
	gl.Status = GameStatusOngoing
	gl.TurnPlayer = 1
	// immediately save to db
	gl.SaveToDB()
}

func (gl *GameLoader) checkGameEnd() {
	if gl.Status != GameStatusOngoing {
		return
	}
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
		logger.GetLogger().Debug("loader: game ended", zap.Uint64("game_id", gl.ID))
		gl.Status = GameStatusEnded
		gl.TurnPlayer = 0
		// immediately save to db
		gl.SaveToDB()
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
	if gl.Status == GameStatusEnded {
		return message.GameErrorMessage(errMsgGameEnded), false
	}

	if msg.Cmd == message.CmdJoin {
		return gl.handleJoin(msg)
	}

	if gl.Status == GameStatusPicking {
		return message.GameErrorMessage(errMsgGameNotStarted), false
	}

	// only current player can do stuff
	if msg.Sender != gl.GameUsers[gl.TurnPlayer-1].UserID {
		return message.GameErrorMessage(errMsgNotYetTurn), false
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
