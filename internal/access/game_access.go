package access

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"gitlab.com/beewar/beewar-be/internal/access/model"
	"gitlab.com/beewar/beewar-be/internal/logger"
	"go.uber.org/zap"
)

/*
CreateGameFromMap initializes a game from a map model, and returns the id.

Provide the following fields (from map model):
 - id
 - type
 - height
 - width
 - player_count
 - terrain_info
 - unit_info

If password is provided, the game will be private (password-protected). Otherwise it will be public.
*/
func CreateGameFromMap(mapModel *model.Map, name, password string, creatorUserID uint64) (uint64, error) {
	const stmtCreateGame = `INSERT INTO game_tab
(type, height, width, player_count, terrain_info, unit_info, map_id, name, password, creator_user_id, time_created, time_modified)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())`

	res, err := db.Exec(stmtCreateGame,
		mapModel.Type, mapModel.Height, mapModel.Width, mapModel.PlayerCount, mapModel.TerrainInfo, mapModel.UnitInfo, mapModel.ID,
		name, password, creatorUserID)
	if err != nil {
		logger.GetLogger().Error("db: insert error", zap.String("table", "game_tab"), zap.Error(err))
		return 0, err
	}
	lastInsertID, err := res.LastInsertId()
	return uint64(lastInsertID), err
}

/*
UpdateGameUsingTx saves a game model to db. If given transaction is nil, db will be used directly.

only updates updatable fields:
 - unit_info
 - status
 - turn_count
 - turn_player
*/
func UpdateGameUsingTx(tx *sql.Tx, game *model.Game) error {
	const stmtUpdateGame = `UPDATE game_tab
SET unit_info=?, status=?, turn_count=?, turn_player=?, time_modified=UNIX_TIMESTAMP()
WHERE id=?`

	var err error
	if tx == nil {
		_, err = db.Exec(stmtUpdateGame,
			game.UnitInfo, game.Status, game.TurnCount, game.TurnPlayer,
			game.ID)
	} else {
		_, err = tx.Exec(stmtUpdateGame,
			game.UnitInfo, game.Status, game.TurnCount, game.TurnPlayer,
			game.ID)
	}
	if err != nil {
		logger.GetLogger().Error("db: update error", zap.String("table", "game_tab"), zap.Error(err))
		return err
	}
	return nil
}

/*
UpdateGameAndGameUser saves a game model, and the game users related to it to db.

only updates updatable fields (game):
 - unit_info
 - status
 - turn_count
 - turn_player
only updates updatable fields (game user):
 - final_rank
 - final_turns
*/
func UpdateGameAndGameUser(game *model.Game, gameUsers []*model.GameUser) error {
	return ExecWithTransaction(func(tx *sql.Tx) error {
		if err := UpdateGameUsingTx(tx, game); err != nil {
			return err
		}
		for _, gu := range gameUsers {
			if err := UpdateGameUserUsingTx(tx, gu); err != nil {
				return err
			}
		}
		return nil
	})
}

// QueryGameByID gets a game from its id
func QueryGameByID(gameID uint64) (*model.Game, error) {
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
		&game.Name,
		&game.Password,
		&game.CreatorUserID,
		&game.Status,
		&game.TurnCount,
		&game.TurnPlayer,
		&game.TimeCreated,
		&game.TimeModified)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.GetLogger().Error("db: query error", zap.String("table", "game_tab"), zap.Error(err))
			return nil, err
		}
		return nil, nil
	}

	return game, nil
}

// QueryGamesByID gets a list of games by id.
// The returned slice will always have the same length as the given game ids, and games are placed in the order given.
func QueryGamesByID(gameIDs []uint64) ([]*model.Game, error) {
	if len(gameIDs) == 0 {
		return make([]*model.Game, 0), nil
	}
	stmt, args, err := sqlx.In(`SELECT * FROM game_tab WHERE id IN (?)`, gameIDs)
	if err != nil {
		logger.GetLogger().Error("db: build sqlx query error", zap.String("table", "game_tab"), zap.Error(err))
		return nil, err
	}
	rows, err := db.Query(stmt, args...)
	if err != nil {
		logger.GetLogger().Error("db: query error", zap.String("table", "game_tab"), zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	games := make([]*model.Game, len(gameIDs))
	for rows.Next() {
		game := &model.Game{}
		err := rows.Scan(
			&game.ID,
			&game.Type,
			&game.Height,
			&game.Width,
			&game.PlayerCount,
			&game.TerrainInfo,
			&game.UnitInfo,
			&game.MapID,
			&game.Name,
			&game.Password,
			&game.CreatorUserID,
			&game.Status,
			&game.TurnCount,
			&game.TurnPlayer,
			&game.TimeCreated,
			&game.TimeModified)
		if err != nil {
			logger.GetLogger().Error("db: query error", zap.String("table", "game_tab"), zap.Error(err))
			return nil, err
		}
		// find index to insert
		for i, id := range gameIDs {
			if id == game.ID {
				games[i] = game
				break
			}
		}
	}
	if err := rows.Err(); err != nil {
		logger.GetLogger().Error("db: query error", zap.String("table", "game_tab"), zap.Error(err))
		return nil, err
	}
	return games, nil
}

// QueryWaitingGames gets all games that are currently not yet started (status = 0)
func QueryWaitingGames() ([]*model.Game, error) {
	rows, err := db.Query(`SELECT * FROM game_tab WHERE status=0 ORDER BY time_created`)
	if err != nil {
		logger.GetLogger().Error("db: query error", zap.String("table", "game_tab"), zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	games := make([]*model.Game, 0)
	for rows.Next() {
		game := &model.Game{}
		err := rows.Scan(
			&game.ID,
			&game.Type,
			&game.Height,
			&game.Width,
			&game.PlayerCount,
			&game.TerrainInfo,
			&game.UnitInfo,
			&game.MapID,
			&game.Name,
			&game.Password,
			&game.CreatorUserID,
			&game.Status,
			&game.TurnCount,
			&game.TurnPlayer,
			&game.TimeCreated,
			&game.TimeModified)
		if err != nil {
			logger.GetLogger().Error("db: query error", zap.String("table", "game_tab"), zap.Error(err))
			return nil, err
		}
		games = append(games, game)
	}
	if err := rows.Err(); err != nil {
		logger.GetLogger().Error("db: query error", zap.String("table", "game_tab"), zap.Error(err))
	}
	return games, nil
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
