package main

import (
	"fmt"
	"gitlab.com/otqee/otqee-be/internal/access"
	"gitlab.com/otqee/otqee-be/internal/logger"
)

func main() {
	logger.InitLogger()
	access.InitAccess()

	user := access.GetUserByUsername("malfple")
	fmt.Println(user)

	access.ShutdownAccess()
	logger.ShutdownLogger()
}
