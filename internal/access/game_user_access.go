package access

import (
	"database/sql"
	"gitlab.com/beewar/beewar-be/internal/access/model"
	"gitlab.com/beewar/beewar-be/internal/logger"
	"go.uber.org/zap"
)

// CreateGameUser creates a game - user link.
func CreateGameUser(gameID, userID uint64, playerOrder uint8) error {
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

/*
UpdateGameUserUsingTx saves a gameUser model to db. If given transaction is nil, db will be used directly.

only updates updatable fields:
 - final_rank
 - final_turns
 - moves_made
*/
func UpdateGameUserUsingTx(tx *sql.Tx, gameUser *model.GameUser) error {
	const stmtUpdateGameUser = `UPDATE game_user_tab
SET final_rank=?, final_turns=?, moves_made=?
WHERE id=?`

	var err error
	if tx == nil {
		_, err = db.Exec(stmtUpdateGameUser,
			gameUser.FinalRank, gameUser.FinalTurns, gameUser.MovesMade,
			gameUser.ID)
	} else {
		_, err = tx.Exec(stmtUpdateGameUser,
			gameUser.FinalRank, gameUser.FinalTurns, gameUser.MovesMade,
			gameUser.ID)
	}
	if err != nil {
		logger.GetLogger().Error("db: update error", zap.String("table", "game_user_tab"), zap.Error(err))
		return err
	}
	return nil
}

// QueryGameUsersByGameID return the game-user link information of a gameID
func QueryGameUsersByGameID(gameID uint64) ([]*model.GameUser, error) {
	rows, err := db.Query(`SELECT * FROM game_user_tab WHERE game_id=? ORDER BY player_order`, gameID)
	if err != nil {
		logger.GetLogger().Error("db: query error", zap.String("table", "game_user_tab"), zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	res := make([]*model.GameUser, 0)
	for rows.Next() {
		gameUser := &model.GameUser{}
		err := rows.Scan(
			&gameUser.ID,
			&gameUser.GameID,
			&gameUser.UserID,
			&gameUser.PlayerOrder,
			&gameUser.FinalRank,
			&gameUser.FinalTurns,
			&gameUser.MovesMade)
		if err != nil {
			logger.GetLogger().Error("db: query error", zap.String("table", "game_user_tab"), zap.Error(err))
			return nil, err
		}
		res = append(res, gameUser)
	}
	if err := rows.Err(); err != nil {
		logger.GetLogger().Error("db: query error", zap.String("table", "game_user_tab"), zap.Error(err))
		return nil, err
	}
	return res, nil
}

// QueryGameUsersByUserID return the game-user link information of a userID
func QueryGameUsersByUserID(userID uint64) ([]*model.GameUser, error) {
	rows, err := db.Query(`SELECT * FROM game_user_tab WHERE user_id=?`, userID)
	if err != nil {
		logger.GetLogger().Error("db: query error", zap.String("table", "game_user_tab"), zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	res := make([]*model.GameUser, 0)
	for rows.Next() {
		gameUser := &model.GameUser{}
		err := rows.Scan(
			&gameUser.ID,
			&gameUser.GameID,
			&gameUser.UserID,
			&gameUser.PlayerOrder,
			&gameUser.FinalRank,
			&gameUser.FinalTurns,
			&gameUser.MovesMade)
		if err != nil {
			logger.GetLogger().Error("db: query error", zap.String("table", "game_user_tab"), zap.Error(err))
			return nil, err
		}
		res = append(res, gameUser)
	}
	if err := rows.Err(); err != nil {
		logger.GetLogger().Error("db: query error", zap.String("table", "game_user_tab"), zap.Error(err))
		return nil, err
	}
	return res, nil
}

// QueryGameUser queries game user by (game_id, user_id)
func QueryGameUser(gameID, userID uint64) (*model.GameUser, error) {
	row := db.QueryRow(`SELECT * FROM game_user_tab WHERE user_id=? AND game_id=? LIMIT 1`, userID, gameID)

	gameUser := &model.GameUser{}
	err := row.Scan(
		&gameUser.ID,
		&gameUser.GameID,
		&gameUser.UserID,
		&gameUser.PlayerOrder,
		&gameUser.FinalRank,
		&gameUser.FinalTurns,
		&gameUser.MovesMade)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.GetLogger().Error("db: query error", zap.String("table", "game_user_tab"), zap.Error(err))
			return nil, err
		}
		return nil, nil
	}

	return gameUser, nil
}

// IsExistGameUserByLink checks if a game user exists by (user_id, game_id) link,
// which means check if a user is registered to a game.
func IsExistGameUserByLink(userID, gameID uint64) bool {
	row := db.QueryRow(`SELECT 1 FROM game_user_tab WHERE user_id=? AND game_id=? LIMIT 1`, userID, gameID)

	var temp int
	err := row.Scan(&temp)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.GetLogger().Error("db: query error", zap.String("table", "game_user_tab"), zap.Error(err))
		}
		return false
	}
	return true
}

// IsExistGameUserByPlayerOrder checks if a game user exists by (game_id, player_order) key,
// which means check if a certain player order in a game is already taken or not.
func IsExistGameUserByPlayerOrder(gameID uint64, playerOrder uint8) bool {
	row := db.QueryRow(`SELECT 1 FROM game_user_tab WHERE game_id=? AND player_order=? LIMIT 1`, gameID, playerOrder)

	var temp int
	err := row.Scan(&temp)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.GetLogger().Error("db: query error", zap.String("table", "game_user_tab"), zap.Error(err))
		}
		return false
	}
	return true
}
