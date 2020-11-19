package access

import (
	"database/sql"
	"errors"
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

func linkGameToUser(gameID, userID int64) error {
	const stmtLinkGameToUser = `INSERT INTO game_user_tab
(game_id, user_id)
VALUES (?, ?)`

	_, err := db.Exec(stmtLinkGameToUser,
		gameID, userID)
	if err != nil {
		logger.GetLogger().Error("db: insert error", zap.String("table", "game_user_tab"), zap.Error(err))
		return err
	}
	return nil
}

// CreateGameFromMap initializes a game from a map, and returns the id
func CreateGameFromMap(mapID int64, userIDs []int64) (int64, error) {
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
(type, width, height, player_count, terrain_info, unit_info, map_id, time_created, time_modified)
VALUES (?, ?, ?, ?, ?, ?, ?, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())`

	res, err := db.Exec(stmtCreateGameFromMap,
		mapp.Type, mapp.Width, mapp.Height, mapp.PlayerCount, mapp.TerrainInfo, mapp.UnitInfo, mapp.ID)
	if err != nil {
		logger.GetLogger().Error("db: insert error", zap.String("table", "game_tab"), zap.Error(err))
		return 0, err
	}
	gameID, err := res.LastInsertId()
	if err != nil {
		return gameID, err
	}

	// link game to users
	for _, userID := range userIDs {
		err := linkGameToUser(gameID, userID)
		if err != nil {
			return 0, err
		}
	}

	return gameID, nil
}

// QueryGameByID gets a game from its id
func QueryGameByID(gameID int64) *model.Game {
	row := db.QueryRow(`SELECT * FROM game_tab WHERE id=? LIMIT 1`, gameID)

	game := &model.Game{}
	err := row.Scan(
		&game.ID,
		&game.Type,
		&game.Width,
		&game.Height,
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