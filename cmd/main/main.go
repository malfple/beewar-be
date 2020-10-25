package main

import (
	"context"
	"gitlab.com/otqee/otqee-be/internal/handler"
	"gitlab.com/otqee/otqee-be/internal/handler/oauth2"
	"gitlab.com/otqee/otqee-be/internal/logger"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

func main() {
	port := int64(3001)
	logger.InitLogger()

	server := &http.Server{
		Addr:         ":" + strconv.FormatInt(port, 10),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      handler.RootRouter(),
	}

	// TODO: move to its own module (not inside handler)
	oauth2.InitOAuth2Server()

	go func() {
		logger.Logger.Info("starting server...", zap.Int64("port", port))
		if err := server.ListenAndServe(); err != nil {
			logger.Logger.Info("server shutting down...", zap.Error(err))
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
	logger.Logger.Info("shutting down...")
	logger.ShutdownLogger()
	os.Exit(0)
}
