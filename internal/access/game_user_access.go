package access

import (
	"gitlab.com/beewar/beewar-be/internal/access/model"
	"gitlab.com/beewar/beewar-be/internal/logger"
	"go.uber.org/zap"
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

// QueryUsersLinkedToGame return the game-user link information of a gameID
func QueryUsersLinkedToGame(gameID uint64) []*model.GameUser {
	rows, err := db.Query(`SELECT * FROM game_user_tab WHERE game_id=?`, gameID)
	if err != nil {
		logger.GetLogger().Error("db: query error", zap.String("table", "game_user_tab"), zap.Error(err))
		return nil
	}
	defer rows.Close()

	var res []*model.GameUser
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

// QueryGamesLinkedToUser return the game-user link information of a userID
func QueryGamesLinkedToUser(userID uint64) []*model.GameUser {
	rows, err := db.Query(`SELECT * FROM game_user_tab WHERE user_id=?`, userID)
	if err != nil {
		logger.GetLogger().Error("db: query error", zap.String("table", "game_user_tab"), zap.Error(err))
		return nil
	}
	defer rows.Close()

	var res []*model.GameUser
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
