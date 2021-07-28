package seeder

import "gitlab.com/beewar/beewar-be/internal/controller/auth"

// SeedUsers inserts default users
func SeedUsers() bool {
	_ = auth.Register("beebot", "beebot", "beebotbeebot")
	_ = auth.Register("malfple@user.com", "malfple", "malfplesecret")
	_ = auth.Register("rapel@user.com", "rapel", "rapelsecret")
	_ = auth.Register("sebas@user.com", "sebas", "sebassecret")
	_ = auth.Register("kyon@user.com", "kyon", "kyonsecret")
	return true
}
