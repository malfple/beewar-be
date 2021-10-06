package seeder

import (
	"gitlab.com/beewar/beewar-be/internal/access"
	"gitlab.com/beewar/beewar-be/internal/controller/auth"
)

// SeedUsers inserts default users
func SeedUsers() bool {
	if user, _ := access.QueryUserByUsername("beebot"); user == nil {
		_ = auth.Register("beebot", "beebot", "beebotbeebot")
	}
	if user, _ := access.QueryUserByUsername("test1"); user == nil {
		_ = auth.Register("test1@user.com", "test1", "testtest")
	}
	if user, _ := access.QueryUserByUsername("test2"); user == nil {
		_ = auth.Register("test2@user.com", "test2", "testtest")
	}
	if user, _ := access.QueryUserByUsername("test3"); user == nil {
		_ = auth.Register("test3@user.com", "test3", "testtest")
	}
	if user, _ := access.QueryUserByUsername("test4"); user == nil {
		_ = auth.Register("test4@user.com", "test4", "testtest")
	}
	return true
}
