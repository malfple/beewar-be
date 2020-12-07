package access

import (
	"database/sql"
	"errors"
	"gitlab.com/otqee/otqee-be/internal/access/formatter"
	"gitlab.com/otqee/otqee-be/internal/access/model"
	"gitlab.com/otqee/otqee-be/internal/logger"
	"go.uber.org/zap"
)

var (
	// ErrMapDoesNotExist is returned when the mapID used for the game is not valid
	ErrMapDoesNotExist = errors.New("map does not exist")
	// ErrUserCountMismatch is returned when the userIDs given doesn't match map's player count
	ErrUserCountMismatch = errors.New("user id given has to be as many as map's player count")
	// ErrUserDoesNotExist is returned when the given userID is invalid
	ErrUserDoesNotExist = errors.New("user does not exist")
)

func linkGameToUser(gameID, userID uint64, playerOrder uint8) error {
	// player order defines the player number of user `userID` in game `gameID`
	const stmtLinkGameToUser = `INSERT INTO game_user_tab
(game_id, user_id, player_order)
VALUES (?, ?, ?)`

	_, err := db.Exec(stmtLinkGameToUser,
		gameID, userID, playerOrder)
	if err != nil {
		logger.GetLogger().Error("db: insert error", zap.String("table", "game_user_tab"), zap.Error(err))
		return err
	}
	return nil
}

// CreateGameFromMap initializes a game from a map, and returns the id
// userIDs should be ordered by the player number.
// userIDs[0] will be player number 1, userIDs[1] player num 2, and so on.
func CreateGameFromMap(mapID uint64, userIDs []uint64) (uint64, error) {
	mapp := QueryMapByID(mapID)
	if mapp == nil {
		return 0, ErrMapDoesNotExist
	}
	if len(userIDs) != int(mapp.PlayerCount) {
		return 0, ErrUserCountMismatch
	}
	for _, userID := range userIDs {
		if !IsExistUserByID(userID) {
			return 0, ErrUserDoesNotExist
		}
	}

	const stmtCreateGameFromMap = `INSERT INTO game_tab
(type, height, width, player_count, terrain_info, unit_info, map_id, time_created, time_modified)
VALUES (?, ?, ?, ?, ?, ?, ?, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())`

	res, err := db.Exec(stmtCreateGameFromMap,
		mapp.Type, mapp.Height, mapp.Width, mapp.PlayerCount, mapp.TerrainInfo, mapp.UnitInfo, mapp.ID)
	if err != nil {
		logger.GetLogger().Error("db: insert error", zap.String("table", "game_tab"), zap.Error(err))
		return 0, err
	}
	gameIDSigned, err := res.LastInsertId()
	gameID := uint64(gameIDSigned)
	if err != nil {
		return gameID, err
	}

	// link game to users
	for i, userID := range userIDs {
		err := linkGameToUser(gameID, userID, uint8(i+1))
		if err != nil {
			return 0, err
		}
	}

	return gameID, nil
}

// UpdateGame saves a game model to db.
// only updates updatable fields
func UpdateGame(game *model.Game) error {
	if err := formatter.ValidateUnitInfo(game.Height, game.Width, game.UnitInfo); err != nil {
		return err
	}

	const stmtUpdateGame = `UPDATE game_tab
SET unit_info=?, turn_count=?, turn_player=?, time_modified=UNIX_TIMESTAMP()
WHERE id=?`

	_, err := db.Exec(stmtUpdateGame,
		game.UnitInfo, game.TurnCount, game.TurnPlayer,
		game.ID)
	if err != nil {
		logger.GetLogger().Error("db: update error", zap.String("table", "game_tab"), zap.Error(err))
		return err
	}
	return nil
}

// QueryGameByID gets a game from its id
func QueryGameByID(gameID uint64) *model.Game {
	row := db.QueryRow(`SELECT * FROM game_tab WHERE id=? LIMIT 1`, gameID)

	game := &model.Game{}
	err := row.Scan(
		&game.ID,
		&game.Type,
		&game.Height,
		&game.Width,
		&game.PlayerCount,
		&game.TerrainInfo,
		&game.UnitInfo,
		&game.MapID,
		&game.TurnCount,
		&game.TurnPlayer,
		&game.TimeCreated,
		&game.TimeModified)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.GetLogger().Error("db: query error", zap.String("table", "game_tab"), zap.Error(err))
		}
		return nil
	}

	return game
}

// IsExistGameByID checks for gameID existence
func IsExistGameByID(gameID uint64) bool {
	row := db.QueryRow(`SELECT 1 FROM game_tab WHERE id=? LIMIT 1`, gameID)

	var temp int
	err := row.Scan(&temp)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.GetLogger().Error("db: query error", zap.String("table", "game_tab"), zap.Error(err))
		}
		return false
	}
	return true
}
