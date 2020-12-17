package user

import "github.com/gorilla/mux"

// RegisterUserRouter builds router for user related stuff
func RegisterUserRouter(router *mux.Router) {
	router.HandleFunc("/get_by_username", HandleUserGetByUsername).Methods("GET")
}

// ReducedUser is model.User but without password
type ReducedUser struct {
	ID          uint64 `json:"id"`
	Email       string `json:"email"`
	Username    string `json:"username"`
	Rating      uint16 `json:"rating"`
	MovesMade   uint64 `json:"moves_made"`
	GamesPlayed uint32 `json:"games_played"`
	TimeCreated int64  `json:"time_created"`
}
