package logger

import (
	"fmt"
	"go.uber.org/zap"
)

// Logger is the default logger
var Logger *zap.Logger

// InitLogger initializes logger
func InitLogger() {
	var err error
	Logger, err = zap.NewDevelopment()
	if err != nil {
		fmt.Println("failed to init logger!")
		fmt.Println(err)
		return
	}

	defer Logger.Sync()
	Logger.Info("init logger")
}

// ShutdownLogger cleans the logger before exiting
func ShutdownLogger() {
	Logger.Info("sync logger to shutdown")
	_ = Logger.Sync()
}
