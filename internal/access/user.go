package access

import (
	"database/sql"
	"gitlab.com/otqee/otqee-be/internal/access/model"
	"gitlab.com/otqee/otqee-be/internal/logger"
	"go.uber.org/zap"
)

// GetUserByUsername gets a single user by username
func GetUserByUsername(username string) *model.User {
	row := db.QueryRow(`SELECT * FROM user_tab WHERE username=?`, username)

	user := &model.User{}
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.TimeCreated)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.GetLogger().Error("db: query error", zap.Error(err))
		}
		return nil
	}
	return user
}
