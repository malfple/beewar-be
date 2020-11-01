package main

import (
	"fmt"
	"gitlab.com/otqee/otqee-be/internal/access"
	"gitlab.com/otqee/otqee-be/internal/logger"
)

func main() {
	logger.InitLogger()
	access.InitAccess()

	//user := access.GetUserByUsername("malfple")
	//fmt.Println(user)

	//mapID, err := access.CreateEmptyMap(0, 2, 3, 1)
	//fmt.Println(mapID, err)

	mapp := access.GetMapByID(2)
	fmt.Println(mapp)

	access.ShutdownAccess()
	logger.ShutdownLogger()
}
