package handler

import (
	"fmt"
	"net/http"
)

// HandlePing just returns pong
func HandlePing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong!")
}
