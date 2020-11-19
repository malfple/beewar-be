package main

import (
	"fmt"
	"gitlab.com/otqee/otqee-be/internal/access"
	"gitlab.com/otqee/otqee-be/internal/auth"
	"gitlab.com/otqee/otqee-be/internal/logger"
)

func main() {
	logger.InitLogger()
	access.InitAccess()

	//user := access.QueryUserByUsername("malfple")
	//fmt.Println(user)

	//mapID, err := access.CreateEmptyMap(0, 2, 3, 1)
	//fmt.Println(mapID, err)

	//terrain1 := make([]byte, 100)
	//for i := 0; i < 10; i++ {
	//	terrain1[rand.Int()%100] = 1
	//}
	//_ = access.UpdateMap(1, 0, 10, 10, "some updated seeded map", terrain1, make([]byte, 0))

	mapp := access.QueryMapByID(1)
	fmt.Println(mapp)

	token := auth.GenerateJWT(123, "some_username")
	fmt.Println(token)

	//gameID, err := access.CreateGameFromMap(1, []int64{1, 2})
	//fmt.Printf("game id: %d, err: %v\n", gameID, err)

	game := access.QueryGameByID(1)
	fmt.Println(game)

	access.ShutdownAccess()
	logger.ShutdownLogger()
}
