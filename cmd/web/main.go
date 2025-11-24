package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/ekaputra07/go-retro/internal/board"
	"github.com/ekaputra07/go-retro/internal/natsutil"
	"github.com/ekaputra07/go-retro/internal/store"
	"github.com/ekaputra07/go-retro/internal/store/natstore"
	"github.com/gorilla/sessions"
)

type app struct {
	config  config
	logger  *slog.Logger
	store   *store.GlobalStore
	manager *board.BoardManager
	session *sessions.CookieStore
	nats    *natsutil.NATS
}

func main() {
	// config
	c := parseConfig()

	// logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// session store
	session := sessions.NewCookieStore([]byte(c.secret))
	session.Options = &sessions.Options{Secure: c.secure}

	// NATS
	natscon := natsutil.Connect(c.natsUrl, c.natsCreds)
	defer natscon.Close()

	// database
	ctx := context.Background()
	db, err := natstore.NewGlobalStore(ctx, natscon, "goretro-global")
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	// db := memstore.NewGlobalStore()

	// board manager
	manager := board.NewBoardManager(logger, natscon, db, strings.Split(c.initialColumns, ","))
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go manager.Start(ctx)

	a := &app{
		config:  c,
		logger:  logger,
		store:   db,
		manager: manager,
		session: session,
		nats:    natscon,
	}

	logger.Info(fmt.Sprintf("%s (%s) running on :%d", appName, appVersion, c.port))
	err = http.ListenAndServe(fmt.Sprintf(":%d", c.port), a.routes())
	logger.Error(err.Error())
	os.Exit(1)
}
