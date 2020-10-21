package main

import (
	"context"
	"fmt"
	"gitlab.com/otqee/otqee-be/internal/view"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

func main() {
	port := int64(3001)

	server := &http.Server{
		Addr:         ":" + strconv.FormatInt(port, 10),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      view.RootRouter(),
	}

	go func() {
		fmt.Printf("Starting server at port %d\n", port)
		if err := server.ListenAndServe(); err != nil {
			fmt.Println(err)
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
	server.Shutdown(ctx)
	fmt.Println("shutting down!")
	os.Exit(0)
}
