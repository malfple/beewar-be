package mymap

import "github.com/gorilla/mux"

// RegisterMapRouter builds router for map related stuff
func RegisterMapRouter(router *mux.Router) {
	router.HandleFunc("/get", HandleMapGet).Methods("GET")
	router.HandleFunc("/list", HandleMapList).Methods("GET")
	router.HandleFunc("/create", HandleMapCreate).Methods("POST")
}
