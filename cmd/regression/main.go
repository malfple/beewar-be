package main

import (
	"gitlab.com/otqee/otqee-be/internal/access"
	"gitlab.com/otqee/otqee-be/internal/logger"
	"gitlab.com/otqee/otqee-be/internal/regression"
)

func main() {
	logger.InitLogger()
	access.InitAccess()

	logger.GetLogger().Info("start regression test")
	if regression.RunRegressionTests() {
		logger.GetLogger().Info("regression test finished successfully")
	} else {
		logger.GetLogger().Info("regression test failed")
	}

	access.ShutdownAccess()
	logger.ShutdownLogger()
}
