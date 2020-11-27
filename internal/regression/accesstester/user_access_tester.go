package accesstester

import (
	"gitlab.com/otqee/otqee-be/internal/access"
	"gitlab.com/otqee/otqee-be/internal/logger"
)

// TestUserAccess runs regression tests for user access
func TestUserAccess() bool {
	logger.GetLogger().Info("user access tester")

	username := "some_unique_username"
	if err := access.CreateUser(username+"@somemail.com", username, "password"); err != nil {
		return false
	}
	// test user query
	user := access.QueryUserByUsername(username)
	if user == nil {
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
