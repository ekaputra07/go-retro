package main

import (
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/ekaputra07/go-retro/internal/board"
	"github.com/ekaputra07/go-retro/internal/models"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	SESSION_NAME = "goretro_session"
	// avatarsCount is the total number of avatars available to choose from.
	// see: web/public/avatars
	AVATARS_COUNT = 12
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func init() {
	gob.Register(uuid.UUID{})
}

func (a *app) health(w http.ResponseWriter, r *http.Request) {
	if a.manager.Healthy() {
		fmt.Fprint(w, "ok")
	} else {
		http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
	}
}

func (a *app) generateBoardID(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, fmt.Sprintf("/b/%s", uuid.New()), http.StatusSeeOther)
}

func (a *app) board(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	session, _ := a.session.Get(r, SESSION_NAME)
	createUser := false

	// create new user if:
	// - session is newdb
	// - user_id in existing user no longer exists
	if session.IsNew {
		createUser = true
	} else {
		userID := session.Values["user_id"].(uuid.UUID)
		if _, err := a.store.Users.Get(ctx, userID); err != nil {
			createUser = true
		}
	}

	if createUser {
		avatarID := rand.Intn(AVATARS_COUNT-1) + 1
		user := models.NewUser(avatarID)
		err := a.store.Users.Create(ctx, user)
		if err != nil {
			a.serverError(w, r, err)
			return
		}

		session.Values["user_id"] = user.ID
		if err := session.Save(r, w); err != nil {
			a.serverError(w, r, err)
			return
		}
		a.logger.Info("new user created", "id", user.ID)
	}

	data := newTemplateData(a.config)
	a.render(w, r, http.StatusOK, data)
}

func (a *app) websocket(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// validate session (make sure user is present) before upgrading the connection
	// TODO: Move this check to a middleware?
	session, _ := a.session.Get(r, SESSION_NAME)
	userID, ok := session.Values["user_id"]
	if !ok {
		a.clientError(w, r, http.StatusUnauthorized, errors.New("session missing user_id"))
		return
	}

	user, err := a.store.Users.Get(ctx, userID.(uuid.UUID))
	if err != nil {
		a.clientError(w, r, http.StatusUnauthorized, fmt.Errorf("error a.store.Users.Get: %s", err.Error()))
		return
	}

	// all good, allow connection
	boardID := r.PathValue("board")
	username := r.URL.Query().Get("u")

	// update name if different
	if user.Name != username {
		user.Name = username
		if err = a.store.Users.Update(ctx, *user); err != nil {
			a.serverError(w, r, fmt.Errorf("error a.store.Users.Update: %s", err.Error()))
			return
		}
	}

	// upgrade to websocket conn
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		a.serverError(w, r, fmt.Errorf("error upgrader.Upgrade: %s", err.Error()))
		return
	}
	defer conn.Close()

	// start board process
	b, err := a.manager.GetOrCreateBoardProcess(ctx, uuid.MustParse(boardID))
	if err != nil {
		a.serverError(w, r, fmt.Errorf("error a.manager.GetOrCreateBoardProcess: %s", err.Error()))
		return
	}

	// create client and add to board
	client := board.NewClient(conn, user, b)
	b.AddClient(client)

	defer client.Stop()
	defer b.RemoveClient(client)

	client.Start()
}
