package main

import (
	"gitlab.com/beewar/beewar-be/configs"
	"gitlab.com/beewar/beewar-be/internal/access"
	"gitlab.com/beewar/beewar-be/internal/logger"
	"gitlab.com/beewar/beewar-be/internal/regression"
	"os"
)

// this has to be run from the root directory of the project
// so that the relative path to files are correct

// the database for regression test, which is defined in this constant, have to be made manually
const regressionDatabaseName = "beewar_regression"

func prepareAndRun() int {
	logger.InitLogger()
	defer logger.ShutdownLogger()
	configs.InitConfigs()

	// doing something that shouldn't have been done. But only for regression testing
	configs.GetServerConfig().Database.Name = regressionDatabaseName
	if configs.GetServerConfig().Database.Name != regressionDatabaseName {
		logger.GetLogger().Error("error set database name")
		return 1
	}

	access.InitAccess()
	defer access.ShutdownAccess()

	if !regression.RunMigration() {
		return 1
	}

	logger.GetLogger().Info("start regression test")
	if !regression.RunRegressionTests() {
		logger.GetLogger().Info("regression test failed")
		return 1
	}

	logger.GetLogger().Info("regression test finished successfully")
	return 0
}

func main() {
	os.Exit(prepareAndRun())
}
