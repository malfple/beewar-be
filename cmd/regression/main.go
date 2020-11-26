package main

import (
	"gitlab.com/otqee/otqee-be/configs"
	"gitlab.com/otqee/otqee-be/internal/access"
	"gitlab.com/otqee/otqee-be/internal/logger"
	"gitlab.com/otqee/otqee-be/internal/regression"
	"go.uber.org/zap"
	"os"
)

func main() {
	logger.InitLogger()
	defer logger.ShutdownLogger()

	if err := os.Setenv(configs.EnvDatabaseName, "otqee_regression"); err != nil {
		logger.GetLogger().Error("error set env var", zap.String("env_var_name", configs.EnvDatabaseName))
		return
	}

	access.InitAccess()
	defer access.ShutdownAccess()

	// TODO: prepare tables using migration

	logger.GetLogger().Info("start regression test")
	if regression.RunRegressionTests() {
		logger.GetLogger().Info("regression test finished successfully")
	} else {
		logger.GetLogger().Info("regression test failed")
	}
}
