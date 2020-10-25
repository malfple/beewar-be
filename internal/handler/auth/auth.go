package auth

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

// RegisterAuthRouter builds router for auth, which handles authentication and authorization
func RegisterAuthRouter(router *mux.Router) {
	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "login")
	})
}
