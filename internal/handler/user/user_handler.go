package user

import "github.com/gorilla/mux"

// RegisterUserRouter builds router for user related stuff
func RegisterUserRouter(router *mux.Router) {
	router.HandleFunc("/get_by_username", HandleUserGetByUsername).Methods("GET")
	router.HandleFunc("/get_many_by_id", HandleUserGetManyByID).Methods("GET")
	router.HandleFunc("/list", HandleUserList).Methods("GET")
}
