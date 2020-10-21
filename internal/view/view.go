package view

import (
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gitlab.com/otqee/otqee-be/internal/middleware"
)

// This file is the root file for all the views
// Handlers and Middlewares are managed here

// RootRouter returns the highest level router
func RootRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", Ping)
	router.Use(middleware.AccessLogMiddleware)
	router.Use(cors.Default().Handler)
	return router
}
