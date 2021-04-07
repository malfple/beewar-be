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

// ExecWithTransaction is a wrapper function for flexible executions of access functions with transaction.
// Provide an anonymous function which contains your executions. Your executions should use the transaction parameter.
// If the function returns error, the transaction will be rollback-ed and the error will be returned.
func ExecWithTransaction(execFunc func(tx *sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		logger.GetLogger().Error("db: begin transaction error", zap.Error(err))
		return err
	}

	if err = execFunc(tx); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			logger.GetLogger().Error("db: fail to rollback", zap.Error(rollbackErr))
			return rollbackErr
		}
		return err
	}

	if err = tx.Commit(); err != nil {
		logger.GetLogger().Error("db: commit transaction error", zap.Error(err))
		return err
	}

	return nil
}
