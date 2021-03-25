package access

import (
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
UpdateGameUser saves a gameUser model to db

only updates updatable fields:
 - final_rank
 - final_turns
*/
func UpdateGameUser(gameUser *model.GameUser) error {
	const stmtUpdateGameUser = `UPDATE game_user_tab
SET final_rank=?, final_turns=?
WHERE id=?`

	_, err := db.Exec(stmtUpdateGameUser,
		gameUser.FinalRank, gameUser.FinalTurns,
		gameUser.ID)
	if err != nil {
		logger.GetLogger().Error("db: update error", zap.String("table", "game_user_tab"), zap.Error(err))
		return err
	}
	return nil
}

// QueryGameUsersByGameID return the game-user link information of a gameID
func QueryGameUsersByGameID(gameID uint64) []*model.GameUser {
	rows, err := db.Query(`SELECT * FROM game_user_tab WHERE game_id=? ORDER BY player_order`, gameID)
	if err != nil {
		logger.GetLogger().Error("db: query error", zap.String("table", "game_user_tab"), zap.Error(err))
		return nil
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
			&gameUser.FinalTurns)
		if err != nil {
			logger.GetLogger().Error("db: query error", zap.String("table", "game_user_tab"), zap.Error(err))
		} else {
			res = append(res, gameUser)
		}
	}
	if err := rows.Err(); err != nil {
		logger.GetLogger().Error("db: query error", zap.String("table", "game_user_tab"), zap.Error(err))
	}
	return res
}

// QueryGameUsersByUserID return the game-user link information of a userID
func QueryGameUsersByUserID(userID uint64) []*model.GameUser {
	rows, err := db.Query(`SELECT * FROM game_user_tab WHERE user_id=?`, userID)
	if err != nil {
		logger.GetLogger().Error("db: query error", zap.String("table", "game_user_tab"), zap.Error(err))
		return nil
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
			&gameUser.FinalTurns)
		if err != nil {
			logger.GetLogger().Error("db: query error", zap.String("table", "game_user_tab"), zap.Error(err))
		} else {
			res = append(res, gameUser)
		}
	}
	if err := rows.Err(); err != nil {
		logger.GetLogger().Error("db: query error", zap.String("table", "game_user_tab"), zap.Error(err))
	}
	return res
}
