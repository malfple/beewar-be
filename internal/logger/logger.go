package logger

import (
	"fmt"
	"go.uber.org/zap"
)

// logger is the default logger
var logger *zap.Logger

// InitLogger initializes logger
func InitLogger() {
	var err error
	logger, err = zap.NewDevelopment()
	if err != nil {
		fmt.Println("failed to init logger!")
		fmt.Println(err)
		return
	}

	defer func() {
		_ = logger.Sync()
	}()
	logger.Info("init logger")
}

// ShutdownLogger cleans the logger before exiting
func ShutdownLogger() {
	logger.Info("sync logger to shutdown")
	_ = logger.Sync()
}

// GetLogger returns the default logger
func GetLogger() *zap.Logger {
	return logger
}
