package main

import (
	"net/http"
)

func (a *app) routes() http.Handler {
	// static file server with cache middleware (max-age: 3600)
	fileServer := staticCache(http.FileServer(http.Dir(a.config.staticDir)), 3600)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", a.generateBoardID)
	mux.HandleFunc("GET /health", a.health)
	mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))
	mux.HandleFunc("GET /b/{board}", a.board)
	mux.HandleFunc("/b/{board}/ws", a.websocket)

	// apply common headers middleware to all routes
	return a.recoverPanic(a.logRequest(commonHeaders(mux)))
}
