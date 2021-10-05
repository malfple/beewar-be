package seeder

import "gitlab.com/beewar/beewar-be/internal/controller/auth"

// SeedUsers inserts default users
func SeedUsers() bool {
	_ = auth.Register("beebot", "beebot", "beebotbeebot")
	_ = auth.Register("test1@user.com", "test1", "testtest")
	_ = auth.Register("test2@user.com", "test2", "testtest")
	_ = auth.Register("test3@user.com", "test3", "testtest")
	_ = auth.Register("test4@user.com", "test4", "testtest")
	return true
}
