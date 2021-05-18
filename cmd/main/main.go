package main

import (
	"context"
	"gitlab.com/beewar/beewar-be/configs"
	"gitlab.com/beewar/beewar-be/internal/access"
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager"
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager/beebot"
	"gitlab.com/beewar/beewar-be/internal/handler"
	"gitlab.com/beewar/beewar-be/internal/logger"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	logger.InitLogger()
	defer logger.ShutdownLogger()
	configs.InitConfigs()
	access.InitAccess()
	defer access.ShutdownAccess()
	gamemanager.InitGameManager()
	defer gamemanager.ShutdownGameManager()
	beebot.InitBeebotRoutines()

	server := &http.Server{
		Addr:         configs.GetServerAddress(),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      handler.RootRouter(),
	}

	go func() {
		logger.GetLogger().Info("starting server...", zap.String("address", server.Addr))
		if err := server.ListenAndServe(); err != nil {
			logger.GetLogger().Info("server shutting down...", zap.Error(err))
		}
	}()

	// graceful shutdown
	c := make(chan os.Signal, 1)
	// quit when SIGINT (ctrl + c)
	signal.Notify(c, os.Interrupt)

	// block until receive signal
	<-c

	// create deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	// wait for connections to close or until deadline
	_ = server.Shutdown(ctx)
	logger.GetLogger().Info("shutting down...")
}
