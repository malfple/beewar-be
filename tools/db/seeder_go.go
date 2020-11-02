package main

import (
	"fmt"
	"gitlab.com/otqee/otqee-be/internal/access"
	"gitlab.com/otqee/otqee-be/internal/logger"
)

func main() {
	logger.InitLogger()
	access.InitAccess()

	_ = access.CreateUser("malfple@user.com", "malfple", "malfplesecret")
	_ = access.CreateUser("rapel@user.com", "rapel", "rapelsecret")
	_ = access.CreateUser("sebas@user.com", "sebas", "sebassecret")
	_ = access.CreateUser("kyon@user.com", "kyon", "kyonsecret")

	mapID, _ := access.CreateEmptyMap(0, 10, 10, "some seeded map", 1)
	fmt.Printf("create map with id: %d\n", mapID)

	access.ShutdownAccess()
	logger.ShutdownLogger()
}
