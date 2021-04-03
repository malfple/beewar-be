package main

import (
	"gitlab.com/beewar/beewar-be/configs"
	"gitlab.com/beewar/beewar-be/internal/access"
	"gitlab.com/beewar/beewar-be/internal/logger"
	"go.uber.org/zap"
	"io/ioutil"
	"strings"
)

func main() {
	logger.InitLogger()
	defer logger.ShutdownLogger()
	configs.InitConfigs()
	access.InitAccess()
	defer access.ShutdownAccess()

	logger.GetLogger().Info("run db migration to prepare tables")
	// load migration file
	migrationFile, err := ioutil.ReadFile("tools/db/migration.sql")
	if err != nil {
		logger.GetLogger().Error("error open migration file", zap.Error(err))
		return
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
		}
	}

	logger.GetLogger().Info("migration run successfully")
}
