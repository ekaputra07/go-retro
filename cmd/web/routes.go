package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (a *app) routes() *mux.Router {
	// static file server
	fileServer := http.FileServer(http.Dir("./ui/public"))

	router := mux.NewRouter()
	router.Use(a.loggingMiddleware)
	router.Use(a.cacheControlMiddleware)
	router.HandleFunc("/health", a.health)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))
	router.HandleFunc("/b/{board}/ws", a.websocket)
	router.HandleFunc("/b/{board}", a.board)
	router.HandleFunc("/", a.generateBoard)
	return router
}
