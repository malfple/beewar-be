package logger

import (
	"fmt"
	"go.uber.org/zap"
)

// Logger is the default production logger
var Logger *zap.Logger

// InitLogger initializes logger
func InitLogger() {
	var err error
	Logger, err = zap.NewProduction()
	if err != nil {
		fmt.Println("failed to init logger!")
		fmt.Println(err)
		return
	}

	Logger.Info("init logger")
}

// ShutdownLogger cleans the logger before exiting
func ShutdownLogger() {
	Logger.Info("sync logger to shutdown")
	_ = Logger.Sync()
}
