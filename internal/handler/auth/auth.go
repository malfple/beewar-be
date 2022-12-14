package auth

import (
	"github.com/gorilla/mux"
)

// RegisterAuthRouter builds router for auth, which handles authentication and authorization
func RegisterAuthRouter(router *mux.Router) {
	router.HandleFunc("/login", HandleLogin).Methods("POST")
	router.HandleFunc("/token", HandleToken).Methods("POST")
	router.HandleFunc("/logout", HandleLogout).Methods("POST")
	router.HandleFunc("/register", HandleRegister).Methods("POST")
}
