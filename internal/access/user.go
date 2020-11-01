package access

import (
	"database/sql"
	"gitlab.com/otqee/otqee-be/internal/access/model"
	"gitlab.com/otqee/otqee-be/internal/logger"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser creates a new user
func CreateUser(email, username, password string) error {
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

// GetUserByUsername gets a single user by username
func GetUserByUsername(username string) *model.User {
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
