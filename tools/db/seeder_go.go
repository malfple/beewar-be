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
			1, 0, 1, 1, 1, 1, 1, 1, 0, 0,
			1, 0, 1, 1, 1, 1, 1, 1, 1, 0,
			0, 1, 1, 1, 1, 1, 1, 1, 1, 0,
			0, 1, 1, 1, 1, 0, 1, 1, 1, 1,
			1, 1, 1, 0, 0, 0, 0, 1, 1, 1,
			1, 1, 1, 0, 0, 0, 0, 1, 1, 1,
			1, 1, 1, 1, 0, 1, 1, 1, 1, 0,
			0, 1, 1, 1, 1, 1, 1, 1, 1, 0,
			0, 1, 1, 1, 1, 1, 1, 1, 0, 1,
			0, 0, 1, 1, 1, 1, 1, 1, 0, 1,
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
		gameID, _ := access.CreateGameFromMap(1, []uint64{2, 4})
		fmt.Printf("create game with id %d \n", gameID)
	}

	// map 2 and game 2
	if access.QueryMapByID(2) == nil {
		mapID, _ := access.CreateEmptyMap(0, 10, 10, "some seeded map", 1)
		fmt.Printf("create map with id: %d\n", mapID)

		terrainInfo := []byte{
			0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 1, 0, 0,
			0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 0,
			0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0,
			0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0,
			0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0,
			0, 0, 1, 1, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0,
		}
		unitInfo := []byte{
			1, 7, 1, 1, 10, 0,
			1, 6, 1, 3, 10, 0,
			1, 8, 1, 3, 10, 0,
			0, 6, 1, 3, 10, 0,
			0, 7, 1, 3, 10, 0,
			2, 6, 1, 3, 10, 0,
			2, 7, 1, 3, 10, 0,
			4, 11, 2, 1, 10, 0,
			4, 10, 2, 3, 10, 0,
			4, 12, 2, 3, 10, 0,
			3, 11, 2, 3, 10, 0,
			3, 12, 2, 3, 10, 0,
			5, 11, 2, 3, 10, 0,
			5, 12, 2, 3, 10, 0,
			10, 2, 3, 1, 10, 0,
			10, 1, 3, 3, 10, 0,
			10, 3, 3, 3, 10, 0,
			9, 2, 3, 3, 10, 0,
			9, 3, 3, 3, 10, 0,
			11, 2, 3, 3, 10, 0,
			11, 3, 3, 3, 10, 0,
			13, 7, 4, 1, 10, 0,
			13, 6, 4, 3, 10, 0,
			13, 8, 4, 3, 10, 0,
			12, 6, 4, 3, 10, 0,
			12, 7, 4, 3, 10, 0,
			14, 6, 4, 3, 10, 0,
			14, 7, 4, 3, 10, 0,
		}

		_ = access.UpdateMap(2, 0, 15, 15, "cross", 4,
			terrainInfo, unitInfo)
	}
	if access.QueryGameByID(2) == nil {
		gameID, _ := access.CreateGameFromMap(2, []uint64{4, 3, 2, 1})
		fmt.Printf("create game with id %d \n", gameID)
	}

	access.ShutdownAccess()
	logger.ShutdownLogger()
}
