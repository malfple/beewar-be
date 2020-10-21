package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
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

	fmt.Printf("Starting server at port %d\n", port)
	if err := http.ListenAndServe(":" + strconv.FormatInt(port, 10), router); err != nil {
		fmt.Println("fail to start server")
	}
}
