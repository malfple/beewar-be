package main

import (
	"fmt"
	"gitlab.com/beewar/beewar-be/configs"
	"gitlab.com/beewar/beewar-be/internal/access"
	"gitlab.com/beewar/beewar-be/internal/logger"
	"gitlab.com/beewar/beewar-be/internal/regression"
	"gitlab.com/beewar/beewar-be/internal/seeder"
	"os"
)

func prepareAndSeed() int {
	logger.InitLogger()
	defer logger.ShutdownLogger()
	configs.InitConfigs()
	access.InitAccess()
	defer access.ShutdownAccess()

	fmt.Println("(DANGEROUS) run migration? (y/n)")
	var response string
	if _, err := fmt.Scanln(&response); err != nil {
		fmt.Println(err)
		return 1
	}
	if response == "y" {
		if !regression.RunMigration() {
			return 1
		}
	}

	if !seeder.SeedUsers() {
		return 1
	}
	if !seeder.SeedSampleMaps() {
		return 1
	}
	if !seeder.SeedCampaignMaps() {
		return 1
	}

	logger.GetLogger().Info("seeder finished successfully")

	return 0
}

func main() {
	os.Exit(prepareAndSeed())
}
