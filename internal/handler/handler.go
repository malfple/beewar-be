package handler

import (
	"github.com/gorilla/mux"
	"gitlab.com/beewar/beewar-be/internal/handler/auth"
	"gitlab.com/beewar/beewar-be/internal/handler/game"
	_map "gitlab.com/beewar/beewar-be/internal/handler/map"
	"gitlab.com/beewar/beewar-be/internal/handler/user"
	"gitlab.com/beewar/beewar-be/internal/middleware"
)

// This file is the root file for all the views
// Handlers and Middlewares are managed here

// RootRouter returns the highest level router
func RootRouter() *mux.Router {
	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api").Subrouter()

	apiRouter.HandleFunc("/", Ping).Methods("GET")
	auth.RegisterAuthRouter(apiRouter.PathPrefix("/auth").Subrouter())
	user.RegisterUserRouter(apiRouter.PathPrefix("/user").Subrouter())

	_map.RegisterMapRouter(apiRouter.PathPrefix("/map").Subrouter())
	game.RegisterGameRouter(apiRouter.PathPrefix("/game").Subrouter())

	apiRouter.HandleFunc("/server_stats", HandleServerStats).Methods("GET")

	router.Use(middleware.AccessLogMiddleware)

	return router
}
