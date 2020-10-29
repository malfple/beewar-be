package handler

import (
	"github.com/gorilla/mux"
	"gitlab.com/otqee/otqee-be/internal/handler/auth"
	"gitlab.com/otqee/otqee-be/internal/middleware"
)

// This file is the root file for all the views
// Handlers and Middlewares are managed here

// RootRouter returns the highest level router
func RootRouter() *mux.Router {
	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api").Subrouter()

	apiRouter.HandleFunc("/", Ping).Methods("GET")
	auth.RegisterAuthRouter(apiRouter.PathPrefix("/auth").Subrouter())
	apiRouter.HandleFunc("/profile", HandleProfile).Methods("GET")

	router.Use(middleware.AccessLogMiddleware)

	return router
}
