package main

import (
	"fmt"
	"gitlab.com/beewar/beewar-be/configs"
	"gitlab.com/beewar/beewar-be/internal/access"
	"gitlab.com/beewar/beewar-be/internal/logger"
	"gitlab.com/beewar/beewar-be/internal/regression"
	"gitlab.com/beewar/beewar-be/internal/seeder"
)

func main() {
	logger.InitLogger()
	defer logger.ShutdownLogger()
	configs.InitConfigs()
	access.InitAccess()
	defer access.ShutdownAccess()

	fmt.Println("(DANGEROUS) run migration? (y/n)")
	var response string
	if _, err := fmt.Scanln(&response); err != nil {
		fmt.Println(err)
		return
	}
	if response == "y" {
		if !regression.RunMigration() {
			return
		}
	}

	seeder.SeedUsers()

	seeder.SeedSampleMaps()
}
