package main

import (
	"net/http"
)

func (a *app) routes() *http.ServeMux {
	// static file server with cache middleware (max-age: 3600)
	fileServer := a.cacheMiddleware(http.FileServer(http.Dir(a.config.staticDir)), 3600)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", a.generateBoard)
	mux.HandleFunc("GET /health", a.health)
	mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))
	mux.HandleFunc("GET /b/{board}", a.board)
	mux.HandleFunc("/b/{board}/ws", a.websocket)
	return mux
}
