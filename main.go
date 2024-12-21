package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/ekaputra07/go-retro/internal/server"
	"github.com/gorilla/mux"
)

func main() {
	ws := server.NewWSServer()

	r := mux.NewRouter()
	r.HandleFunc("/health", healthHandler)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web"))))
	r.Handle("/b/{board}/ws", ws)
	r.HandleFunc("/b/{board}", boardHandler)
	r.HandleFunc("/", generateBoardHandler)

	// create and start the server
	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Printf("listening on 127.0.0.1:8080...")
	log.Fatal(srv.ListenAndServe())
}

func generateBoardHandler(w http.ResponseWriter, r *http.Request) {
	id := petname.Generate(3, "-")
	http.Redirect(w, r, fmt.Sprintf("/b/%s", id), http.StatusSeeOther)
}

func boardHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/index.html")
	return
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
}
