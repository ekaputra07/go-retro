package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ekaputra07/go-retro/internal/server"
	"github.com/ekaputra07/go-retro/internal/storage"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var (
	host string = os.Getenv("GORETRO_HOST")
	port string = os.Getenv("GORETRO_PORT")
)

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	db := storage.NewMemoryStore()
	ws := server.NewWSServer(db)
	stopWS := make(chan struct{})
	go ws.Start(stopWS)

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

	go func() {
		log.Printf("http-server running on %s...", hostPort)
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http-server not started: %v", err)
		}
		log.Println("http-server shutting down...")
	}()

	<-sigChan
	stopWS <- struct{}{}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("http-server shutdown error: %v", err)
	}
	log.Println("http-server shutdown complete")
}

func loggingMiddleware(next http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, next)
}

func generateBoardHandler(w http.ResponseWriter, r *http.Request) {
	id := uuid.New()
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
