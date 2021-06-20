package access

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"gitlab.com/beewar/beewar-be/internal/access/model"
	"gitlab.com/beewar/beewar-be/internal/logger"
	"go.uber.org/zap"
)

// CreateUser creates a new user. password has to be already hashed
func CreateUser(email, username, password string) error {
	_, err := db.Exec(`INSERT INTO user_tab(email, username, password, time_created) VALUES (?, ?, ?, UNIX_TIMESTAMP())`,
		email, username, password)
	if err != nil {
		logger.GetLogger().Error("db: insert error", zap.String("table", "user_tab"), zap.Error(err))
		return err
	}

	return nil
}

// QueryUserByUsername gets a single user by username
func QueryUserByUsername(username string) (*model.User, error) {
	row := db.QueryRow(`SELECT * FROM user_tab WHERE username=? LIMIT 1`, username)

	user := &model.User{}
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.Rating,
		&user.MovesMade,
		&user.GamesPlayed,
		&user.TimeCreated)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.GetLogger().Error("db: query error", zap.String("table", "user_tab"), zap.Error(err))
			return nil, err
		}
		return nil, nil
	}
	return user, nil
}

// QueryUsersByID gets a list of users by id.
// The returned slice will always have the same length as the given user ids, and users are placed in the order given.
func QueryUsersByID(userIDs []uint64) ([]*model.User, error) {
	if len(userIDs) == 0 {
		return make([]*model.User, 0), nil
	}
	stmt, args, err := sqlx.In(`SELECT * FROM user_tab WHERE id IN (?)`, userIDs)
	if err != nil {
		logger.GetLogger().Error("db: build sqlx query error", zap.String("table", "user_tab"), zap.Error(err))
		return nil, err
	}
	rows, err := db.Query(stmt, args...)
	if err != nil {
		logger.GetLogger().Error("db: query error", zap.String("table", "user_tab"), zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	users := make([]*model.User, len(userIDs))
	for rows.Next() {
		user := &model.User{}
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Username,
			&user.Password,
			&user.Rating,
			&user.MovesMade,
			&user.GamesPlayed,
			&user.TimeCreated)
		if err != nil {
			logger.GetLogger().Error("db: query error", zap.String("table", "user_tab"), zap.Error(err))
			return nil, err
		}
		// find index to insert
		for i, id := range userIDs {
			if id == user.ID {
				users[i] = user
				break
			}
		}
	}
	if err := rows.Err(); err != nil {
		logger.GetLogger().Error("db: query error", zap.String("table", "user_tab"), zap.Error(err))
		return nil, err
	}
	return users, nil
}

// QueryUsers gets a list of users (slow because of limit offset. feel free to optimize)
func QueryUsers(limit, offset int) ([]*model.User, error) {
	rows, err := db.Query(`SELECT * FROM user_tab LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		logger.GetLogger().Error("db: query error", zap.String("table", "user_tab"), zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	users := make([]*model.User, 0)
	for rows.Next() {
		user := &model.User{}
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Username,
			&user.Password,
			&user.Rating,
			&user.MovesMade,
			&user.GamesPlayed,
			&user.TimeCreated)
		if err != nil {
			logger.GetLogger().Error("db: query error", zap.String("table", "user_tab"), zap.Error(err))
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		logger.GetLogger().Error("db: query error", zap.String("table", "user_tab"), zap.Error(err))
		return nil, err
	}
	return users, nil
}

// IsExistUserByID checks for userID existence
func IsExistUserByID(userID uint64) bool {
	row := db.QueryRow(`SELECT 1 FROM user_tab WHERE id=? LIMIT 1`, userID)

	var temp int
	err := row.Scan(&temp)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.GetLogger().Error("db: query error", zap.String("table", "user_tab"), zap.Error(err))
		}
		return false
	}
	return true
}

// DeleteUserByUsername is a DANGEROUS function that deletes user
// returns rows affected
func DeleteUserByUsername(username string) int64 {
	res, err := db.Exec(`DELETE FROM user_tab WHERE username=?`, username)
	if err != nil {
		logger.GetLogger().Error("db: delete error", zap.String("table", "user_tab"), zap.Error(err))
		return 0
	}
	rowsAffected, _ := res.RowsAffected()
	return rowsAffected
}
