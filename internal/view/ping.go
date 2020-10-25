package view

import (
	"fmt"
	"net/http"
)

// Ping just returns pong
func Ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong!")
}
