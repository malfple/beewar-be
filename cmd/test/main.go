package main

import (
	"fmt"
	"gitlab.com/otqee/otqee-be/internal/access"
	"gitlab.com/otqee/otqee-be/internal/gamemanager/loader"
	"gitlab.com/otqee/otqee-be/internal/logger"
)

func main() {
	logger.InitLogger()
	access.InitAccess()

	//mapp := access.QueryMapByID(1)
	//fmt.Println(mapp)
	//
	//token := auth.GenerateJWT(123, "some_username")
	//fmt.Println(token)

	//gameID, err := access.CreateGameFromMap(1, []int64{1, 2})
	//fmt.Printf("game id: %d, err: %v\n", gameID, err)

	game := access.QueryGameByID(1)
	fmt.Println(game)

	gameObj := loader.NewGameLoader(1)
	fmt.Println(gameObj.Terrain)
	for i, row := range gameObj.Units {
		for j, unit := range row {
			if unit == nil {
				continue
			}
			fmt.Printf("unit at %d %d, %T %v\n", i, j, unit, unit)
		}
	}

	access.ShutdownAccess()
	logger.ShutdownLogger()
}
