package main

import (
	"net/http"

	"github.com/ekaputra07/go-retro/web/ui"
)

func (a *app) routes() http.Handler {
	fileServer := http.FileServer(ui.EmbeddedUiFS())

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", a.generateBoardID)
	mux.HandleFunc("GET /health", a.health)
	mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))
	mux.HandleFunc("GET /b/{board}", a.board)
	mux.HandleFunc("/b/{board}/ws", a.websocket)

	// apply common headers middleware to all routes
	return a.recoverPanic(a.logRequest(commonHeaders(mux)))
}
