package accesstester

import (
	"gitlab.com/beewar/beewar-be/internal/access"
	"gitlab.com/beewar/beewar-be/internal/logger"
	"go.uber.org/zap"
)

// TestUserAccess runs regression tests for user access
func TestUserAccess() bool {
	logger.GetLogger().Info("user access tester")

	username := "some_unique_username"
	if err := access.CreateUser(username+"@somemail.com", username, "password"); err != nil {
		return false
	}
	username2 := "username2"
	if err := access.CreateUser(username2+"@somemail.com", username2, "password"); err != nil {
		return false
	}
	// test user query
	user, err := access.QueryUserByUsername(username)
	if user == nil || err != nil {
		return false
	}
	user2, err := access.QueryUserByUsername(username2)
	if user2 == nil || err != nil {
		return false
	}
	users, err := access.QueryUsersByID([]uint64{user2.ID, user.ID})
	if err != nil {
		return false
	}
	if len(users) != 2 {
		logger.GetLogger().Error("expected 2 users", zap.Int("actual", len(users)))
		return false
	}
	users2, err := access.QueryUsers(10, 0)
	if err != nil {
		return false
	}
	if len(users2) != 2 {
		logger.GetLogger().Error("expected 2 users", zap.Int("actual", len(users2)))
		return false
	}
	if users[0].ID != user2.ID || users[1].ID != user.ID {
		logger.GetLogger().Error("wrong order when batch query user")
		return false
	}
	if !access.IsExistUserByID(user.ID) {
		return false
	}
	if access.IsExistUserByID(0) {
		return false
	}
	if rowsAffected := access.DeleteUserByUsername(username); rowsAffected != 1 {
		return false
	}
	return true
}
