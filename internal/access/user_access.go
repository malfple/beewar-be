package access

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"gitlab.com/otqee/otqee-be/internal/access/model"
	"gitlab.com/otqee/otqee-be/internal/logger"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrUsernameLength is returned when username length is out of constraint
	ErrUsernameLength = errors.New("username length must be at least 1 and at most 255")
	// ErrPasswordLength is returned when password length is out of constraint
	ErrPasswordLength = errors.New("password length must be at least 8 and at most 255")
)

const (
	userMaxUsernameLength = 255
	userMinPasswordLength = 8
	userMaxPasswordLength = 255
)

// CreateUser creates a new user
func CreateUser(email, username, password string) error {
	if len(username) < 1 || len(username) > userMaxUsernameLength {
		return ErrUsernameLength
	}
	if len(password) < userMinPasswordLength || len(password) > userMaxPasswordLength {
		return ErrPasswordLength
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.GetLogger().Error("error bcrypt", zap.Error(err))
		return err
	}

	_, err = db.Exec(`INSERT INTO user_tab(email, username, password, time_created) VALUES (?, ?, ?, UNIX_TIMESTAMP())`,
		email, username, passwordHash)
	if err != nil {
		logger.GetLogger().Error("db: insert error", zap.String("table", "user_tab"), zap.Error(err))
		return err
	}

	return nil
}

// QueryUserByUsername gets a single user by username
func QueryUserByUsername(username string) *model.User {
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
		}
		return nil
	}
	return user
}

// QueryUsersByID gets a list of users by id
func QueryUsersByID(userIDs []uint64) []*model.User {
	stmt, args, err := sqlx.In(`SELECT * FROM user_tab WHERE id IN (?)`, userIDs)
	if err != nil {
		logger.GetLogger().Error("db: build sqlx query error", zap.String("table", "user_tab"), zap.Error(err))
		return nil
	}
	rows, err := db.Query(stmt, args...)
	if err != nil {
		logger.GetLogger().Error("db: query error", zap.String("table", "user_tab"), zap.Error(err))
		return nil
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
		} else {
			// find index to insert
			for i, id := range userIDs {
				if id == user.ID {
					users[i] = user
					break
				}
			}
		}
	}
	if err := rows.Err(); err != nil {
		logger.GetLogger().Error("db: query error", zap.String("table", "user_tab"), zap.Error(err))
	}
	return users
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
