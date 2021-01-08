package access

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql" // mysql driver
	"gitlab.com/beewar/beewar-be/configs"
	"gitlab.com/beewar/beewar-be/internal/logger"
	"go.uber.org/zap"
)

// This package contains data access functions (from db)

var db *sql.DB

// OpenDBConnection opens db connection
func OpenDBConnection() *sql.DB {
	logger.GetLogger().Info("db: opening db connection...")
	db, err := sql.Open("mysql", configs.GetDatabaseMySQLDataSourceName())
	if err != nil {
		logger.GetLogger().Error("db: error opening db connection", zap.Error(err))
		return nil
	}
	return db
}

// InitAccess inits the database client
func InitAccess() {
	db = OpenDBConnection()
}

// ShutdownAccess close the database client
func ShutdownAccess() {
	logger.GetLogger().Info("db: closing...")
	err := db.Close()
	if err != nil {
		logger.GetLogger().Error("db: error closing", zap.Error(err))
	}
}

// GetDBClient returns the default sql.DB object
func GetDBClient() *sql.DB {
	return db
}
