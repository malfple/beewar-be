package main

import (
	"fmt"
	"gitlab.com/otqee/otqee-be/internal/access"
	"gitlab.com/otqee/otqee-be/internal/logger"
)

func main() {
	logger.InitLogger()
	access.InitAccess()

	// users
	_ = access.CreateUser("malfple@user.com", "malfple", "malfplesecret")
	_ = access.CreateUser("rapel@user.com", "rapel", "rapelsecret")
	_ = access.CreateUser("sebas@user.com", "sebas", "sebassecret")
	_ = access.CreateUser("kyon@user.com", "kyon", "kyonsecret")

	// map 1
	if access.QueryMapByID(1) == nil {
		mapID, _ := access.CreateEmptyMap(0, 10, 10, "some seeded map", 1)
		fmt.Printf("create map with id: %d\n", mapID)

		terrainInfo := []byte{
			0, 0, 1, 1, 1, 1, 1, 1, 0, 0,
			0, 0, 1, 1, 1, 1, 1, 1, 1, 0,
			0, 1, 1, 1, 1, 1, 1, 1, 1, 0,
			0, 1, 1, 1, 1, 0, 1, 1, 1, 1,
			1, 1, 1, 0, 0, 0, 0, 1, 1, 1,
			1, 1, 1, 0, 0, 0, 0, 1, 1, 1,
			1, 1, 1, 1, 0, 1, 1, 1, 1, 0,
			0, 1, 1, 1, 1, 1, 1, 1, 1, 0,
			0, 1, 1, 1, 1, 1, 1, 1, 0, 0,
			0, 0, 1, 1, 1, 1, 1, 1, 0, 0,
		}
		unitInfo := []byte{
			5, 1, 1, 1, 10, 0,
			4, 1, 1, 3, 10, 0,
			3, 1, 1, 3, 10, 0,
			6, 1, 1, 3, 10, 0,
			7, 1, 1, 3, 10, 0,
			4, 8, 2, 1, 10, 0,
			3, 8, 2, 3, 10, 0,
			2, 8, 2, 3, 10, 0,
			5, 8, 2, 3, 10, 0,
			6, 8, 2, 3, 10, 0,
		}

		_ = access.UpdateMap(1, 0, 10, 10, "some updated seeded map", 2,
			terrainInfo, unitInfo)
	}

	// game 1
	if access.QueryGameByID(1) == nil {
		gameID, _ := access.CreateGameFromMap(1, []uint64{1, 2})
		fmt.Printf("create game with id %d \n", gameID)
	}

	access.ShutdownAccess()
	logger.ShutdownLogger()
}
