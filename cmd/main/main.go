package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong!")
}

func AccessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		fmt.Println("accessed")
	})
}

func main() {
	port := int64(3001)

	router := mux.NewRouter()
	router.HandleFunc("/", Ping)
	router.Use(AccessLogMiddleware)
	router.Use(cors.Default().Handler)

	server := &http.Server{
		Addr: ":" + strconv.FormatInt(port, 10),
		WriteTimeout: time.Second * 15,
		ReadTimeout: time.Second * 15,
		IdleTimeout: time.Second * 60,
		Handler: router,
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 15)
	defer cancel()
	// wait for connections to close or until deadline
	server.Shutdown(ctx)
	fmt.Println("shutting down!")
	os.Exit(0)
}
