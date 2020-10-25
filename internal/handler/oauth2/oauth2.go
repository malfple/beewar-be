package oauth2

import "github.com/gorilla/mux"

// BuildOAuth2Router builds OAuth2 sub-router
func BuildOAuth2Router(router *mux.Router) {
	router.HandleFunc("/authorize", HandleAuthorize)
	router.HandleFunc("/token", HandleToken)
}
