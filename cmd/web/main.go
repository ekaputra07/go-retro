package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/ekaputra07/go-retro/internal/board"
	"github.com/ekaputra07/go-retro/internal/storage"
	"github.com/gorilla/sessions"
)

// config stores all configurable values
type config struct {
	port      int
	staticDir string
	secret    string
	secure    bool
}

type app struct {
	config  config
	logger  *slog.Logger
	db      storage.Storage
	manager *board.BoardManager
	session *sessions.CookieStore
}

func main() {
	// config
	conf := config{}
	flag.IntVar(&conf.port, "port", 8080, "Port to listen")
	flag.StringVar(&conf.staticDir, "staticDir", "./web/public", "Directory of static files")
	flag.StringVar(&conf.secret, "secret", os.Getenv("GORETRO_SESSION_SECRET"), "Session secret")
	flag.BoolVar(&conf.secure, "secure", false, "Secure cookie by default")
	flag.Parse()

	// make sure secret is not empty
	if conf.secret == "" {
		fmt.Println(
			"Secret is missing!",
			"Set secret via environment variable `GORETRO_SESSION_SECRET` (recommended)",
			"or via `-secret` flag.",
		)
		os.Exit(1)
	}

	// logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// session store
	session := sessions.NewCookieStore([]byte(conf.secret))
	session.Options = &sessions.Options{Secure: conf.secure}

	// database
	db := storage.NewMemoryStore()

	// board manager
	// start WS server
	manager := board.NewBoardManager(logger, db)
	ctx, cancel := context.WithCancel(context.Background())
	go manager.Start(ctx)
	defer cancel()

	a := &app{
		config:  conf,
		logger:  logger,
		db:      db,
		manager: manager,
		session: session,
	}

	logger.Info(fmt.Sprintf("Server running on :%d", conf.port))
	err := http.ListenAndServe(fmt.Sprintf(":%d", conf.port), a.routes())
	logger.Error(err.Error())
	os.Exit(1)
}
