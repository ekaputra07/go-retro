package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"github.com/ekaputra07/go-retro/internal/board"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const SESSION_NAME = "goretro_session"

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var boardTpl = template.Must(template.ParseGlob("web/templates/*.html"))

func init() {
	gob.Register(uuid.UUID{})
}

func (a *app) cacheMiddleware(next http.Handler, maxAgeSecond int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAgeSecond))
		next.ServeHTTP(w, r)
	})
}

func (a *app) health(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "ok")
}

func (a *app) generateBoard(w http.ResponseWriter, r *http.Request) {
	id := uuid.New()
	http.Redirect(w, r, fmt.Sprintf("/b/%s", id), http.StatusSeeOther)
}

func (a *app) board(w http.ResponseWriter, r *http.Request) {
	session, _ := a.session.Get(r, SESSION_NAME)
	createUser := false

	// create new user if:
	// - session is new
	// - user_id in existing user no longer exists
	if session.IsNew {
		createUser = true
	} else {
		userID := session.Values["user_id"].(uuid.UUID)
		if _, err := a.db.GetUser(userID); err != nil {
			createUser = true
		}
	}

	if createUser {
		u, err := a.db.CreateUser()
		if err != nil {
			a.serverError(w, r, err)
			return
		}

		session.Values["user_id"] = u.ID
		if err := session.Save(r, w); err != nil {
			a.serverError(w, r, err)
			return
		}
		a.logger.Info("new user created", "id", u.ID)
	}

	if err := boardTpl.ExecuteTemplate(w, "base", nil); err != nil {
		a.serverError(w, r, err)
	}
}

func (a *app) websocket(w http.ResponseWriter, r *http.Request) {
	// validate session (make sure user is present) before upgrading the connection
	// TODO: Move this check to a middleware?
	session, _ := a.session.Get(r, SESSION_NAME)
	userID, ok := session.Values["user_id"]
	if !ok {
		a.clientError(w, r, http.StatusUnauthorized, errors.New("session missing user_id"))
		return
	}

	user, err := a.db.GetUser(userID.(uuid.UUID))
	if err != nil {
		a.clientError(w, r, http.StatusUnauthorized, err)
		return
	}

	// all good, allow connection
	boardID := r.PathValue("board")
	username := r.URL.Query().Get("u")

	// update name if different
	if user.Name != username {
		user.Name = username
		if err = a.db.UpdateUser(user); err != nil {
			a.serverError(w, r, err)
			return
		}
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		a.serverError(w, r, err)
		return
	}
	defer conn.Close()

	// start board process
	b := a.manager.GetOrStartBoard(uuid.MustParse(boardID))
	// create client and add to board
	client := board.NewClient(conn, user, b)
	b.AddClient(client)

	defer client.Stop()
	defer b.RemoveClient(client)

	client.Start()
}
