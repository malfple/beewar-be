package main

import (
	"fmt"
	"gitlab.com/otqee/otqee-be/internal/access"
	"gitlab.com/otqee/otqee-be/internal/logger"
	"math/rand"
)

func main() {
	logger.InitLogger()
	access.InitAccess()

	//user := access.QueryUserByUsername("malfple")
	//fmt.Println(user)

	//mapID, err := access.CreateEmptyMap(0, 2, 3, 1)
	//fmt.Println(mapID, err)

	terrain1 := make([]byte, 100)
	for i := 0; i < 10; i++ {
		terrain1[rand.Int()%100] = 1
	}
	_ = access.UpdateMap(1, 0, 10, 10, "some updated seeded map", terrain1, make([]byte, 0))

	mapp := access.QueryMapByID(1)
	fmt.Println(mapp)

	access.ShutdownAccess()
	logger.ShutdownLogger()
}
