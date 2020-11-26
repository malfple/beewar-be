package accesstester

import "gitlab.com/otqee/otqee-be/internal/access"

// TestUserAccess runs regression tests for user access
func TestUserAccess() bool {
	// TODO: make access function for create and delete
	// test user query
	user := access.QueryUserByUsername("malfple")
	if user == nil {
		return false
	}
	if !access.IsExistUserByID(1) {
		return false
	}
	if access.IsExistUserByID(0) {
		return false
	}
	return true
}
