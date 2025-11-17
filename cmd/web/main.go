package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/ekaputra07/go-retro/internal/board"
	"github.com/ekaputra07/go-retro/internal/nats"
	"github.com/ekaputra07/go-retro/internal/store"
	nstore "github.com/ekaputra07/go-retro/internal/store/nats"
	"github.com/gorilla/sessions"
)

// config stores all configurable values
type config struct {
	port           int
	staticDir      string
	secret         string
	initialColumns string
	secure         bool
	enableTimer    bool
	enableStandup  bool
}

type app struct {
	config  config
	logger  *slog.Logger
	store   *store.GlobalStore
	manager *board.BoardManager
	session *sessions.CookieStore
}

func parseConfig() config {
	conf := config{}
	flag.IntVar(&conf.port, "port", 8080, "Port to listen")
	flag.StringVar(&conf.staticDir, "staticDir", "./web/public", "Directory of static files")
	flag.StringVar(&conf.secret, "secret", os.Getenv("GORETRO_SESSION_SECRET"), "Session secret")
	flag.StringVar(&conf.initialColumns, "initialColumns", "Good,Bad,Questions,Emoji", "Initial board columns")
	flag.BoolVar(&conf.secure, "secure", false, "Secure cookie by default")
	flag.BoolVar(&conf.enableTimer, "timer", true, "Enable timer feature")
	flag.BoolVar(&conf.enableStandup, "standup", true, "Enable standup feature")
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
	return conf
}

func main() {
	ctx := context.Background()

	// config
	c := parseConfig()

	// logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// session store
	session := sessions.NewCookieStore([]byte(c.secret))
	session.Options = &sessions.Options{Secure: c.secure}

	// NATS
	nts := nats.Setup(os.Getenv("NATS_URL"))
	defer nts.Close()

	// database
	db, err := nstore.NewGlobalStore(ctx, nts, "goretro-global")
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// board manager
	manager := board.NewBoardManager(logger, nts, db, strings.Split(c.initialColumns, ","))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go manager.Start(ctx)

	a := &app{
		config:  c,
		logger:  logger,
		store:   db,
		manager: manager,
		session: session,
	}

	logger.Info(fmt.Sprintf("Server running on :%d", c.port))
	err = http.ListenAndServe(fmt.Sprintf(":%d", c.port), a.routes())
	logger.Error(err.Error())
	os.Exit(1)
}
