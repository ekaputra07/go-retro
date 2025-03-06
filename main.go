package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ekaputra07/go-retro/internal/storage"
)

var (
	host = os.Getenv("GORETRO_HOST")
	port = os.Getenv("GORETRO_PORT")
)

func main() {
	a := &app{db: storage.NewMemoryStore()}
	stop := a.start()
	defer stop()

	// create and start the server
	if host == "" {
		host = "127.0.0.1"
	}
	if port == "" {
		port = "8080"
	}
	hostPort := fmt.Sprintf("%s:%s", host, port)
	srv := &http.Server{
		Handler:      a.router,
		Addr:         hostPort,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("web-server running on %s...", hostPort)
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("web-server not started: %v", err)
	}
}
