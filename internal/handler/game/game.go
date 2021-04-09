package game

import (
	"github.com/gorilla/mux"
)

// RegisterGameRouter builds router for game related stuff
func RegisterGameRouter(router *mux.Router) {
	router.HandleFunc("/ws", HandleGameWS)
	router.HandleFunc("/my_games", HandleMyGames).Methods("GET")
}
