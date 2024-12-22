package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/ekaputra07/go-retro/internal/server"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var (
	host string = os.Getenv("GORETRO_HOST")
	port string = os.Getenv("GORETRO_PORT")
)

func main() {
	ws := server.NewWSServer()

	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.HandleFunc("/health", healthHandler)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web"))))
	r.Handle("/b/{board}/ws", ws)
	r.HandleFunc("/b/{board}", boardHandler)
	r.HandleFunc("/", generateBoardHandler)

	// create and start the server
	if host == "" {
		host = "127.0.0.1"
	}
	if port == "" {
		port = "8080"
	}
	hostPort := fmt.Sprintf("%s:%s", host, port)
	srv := &http.Server{
		Handler:      r,
		Addr:         hostPort,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Printf("go-retro running on %s...", hostPort)
	log.Fatal(srv.ListenAndServe())
}

func loggingMiddleware(next http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, next)
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
