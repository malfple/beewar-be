package main

import (
	"gitlab.com/otqee/otqee-be/configs"
	"gitlab.com/otqee/otqee-be/internal/access"
	"gitlab.com/otqee/otqee-be/internal/logger"
	"gitlab.com/otqee/otqee-be/internal/regression"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"strings"
)

// this has to be run from the root directory of the project
// so that the relative path to files are correct

// the database for regression test, which is defined in this constant, have to be made manually
const regressionDatabaseName = "otqee_regression"

func runMigration() bool {
	logger.GetLogger().Info("run db migration to prepare tables")
	// load migration file
	migrationFile, err := ioutil.ReadFile("tools/db/migration.sql")
	if err != nil {
		logger.GetLogger().Error("error open migration file", zap.Error(err))
		return false
	}
	// split statements
	migrationStmts := strings.Split(string(migrationFile), ";\n")
	// exclude the last one, which is empty
	migrationStmts = migrationStmts[:len(migrationStmts)-1]
	// run migration
	for _, stmt := range migrationStmts {
		_, err := access.GetDBClient().Exec(stmt)
		if err != nil {
			logger.GetLogger().Error("error running migration", zap.Error(err))
			return false
		}
	}
	return true
}

func prepareAndRun() int {
	logger.InitLogger()
	defer logger.ShutdownLogger()

	if err := os.Setenv(configs.EnvDatabaseName, regressionDatabaseName); err != nil {
		logger.GetLogger().Error("error set env var", zap.String("env_var_name", configs.EnvDatabaseName))
		return 1
	}

	access.InitAccess()
	defer access.ShutdownAccess()

	if !runMigration() {
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
