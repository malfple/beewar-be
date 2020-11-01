package main

import (
	"gitlab.com/otqee/otqee-be/internal/access"
	"gitlab.com/otqee/otqee-be/internal/logger"
)

func main() {
	logger.InitLogger()
	access.InitAccess()

	_ = access.CreateUser("malfple@user.com", "malfple", "malfple")
	_ = access.CreateUser("rapel@user.com", "rapel", "rapel")
	_ = access.CreateUser("sebas@user.com", "sebas", "sebas")
	_ = access.CreateUser("kyon@user.com", "kyon", "kyon")

	access.ShutdownAccess()
	logger.ShutdownLogger()
}
