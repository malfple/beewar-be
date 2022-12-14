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
	userAgain, err := access.QueryUserByID(user.ID)
	if userAgain == nil || err != nil {
		return false
	}
	if userAgain.Username != user.Username {
		logger.GetLogger().Error("query user by id mismatch")
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
	// update user
	user.MovesMade = 5
	user.GamesPlayed = 1
	if err := access.UpdateUserUsingTx(nil, user); err != nil {
		return false
	}
	userAgain, _ = access.QueryUserByUsername(user.Username)
	if userAgain.MovesMade != 5 || userAgain.GamesPlayed != 1 {
		logger.GetLogger().Error("mismatch update user",
			zap.Uint64("actual moves made", userAgain.MovesMade),
			zap.Uint32("actual games played", userAgain.GamesPlayed))
	}
	// delete
	if rowsAffected := access.DeleteUserByUsername(username); rowsAffected != 1 {
		return false
	}
	return true
}
