package user

import "github.com/gorilla/mux"

// RegisterUserRouter builds router for user related stuff
func RegisterUserRouter(router *mux.Router) {
	router.HandleFunc("/get_by_username", HandleUserGetByUsername).Methods("GET")
}
