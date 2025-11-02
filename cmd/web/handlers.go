package main

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/ekaputra07/go-retro/internal/board"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

const SESSION_NAME = "goretro_session"

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var boardTpl = template.Must(template.ParseGlob("ui/templates/*.html"))

func init() {
	gob.Register(uuid.UUID{})
}

func (a *app) loggingMiddleware(next http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, next)
}

func (a *app) cacheControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/static/") {
			w.Header().Set("Cache-Control", "public, max-age=3600")
		}
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
		a.logger.Info(fmt.Sprintf("new user created with id=%s", u.ID))
	}

	if err := boardTpl.ExecuteTemplate(w, "base", nil); err != nil {
		a.serverError(w, r, err)
	}
}

func (a *app) websocket(w http.ResponseWriter, r *http.Request) {
	// validate session (make sure user is present) before upgrading the connection
	// TODO: Move this check to a middleware
	session, _ := a.session.Get(r, SESSION_NAME)
	userID := session.Values["user_id"].(uuid.UUID)

	user, err := a.db.GetUser(userID)
	if err != nil {
		a.clientError(w, http.StatusUnauthorized)
		return
	}

	// all good, allow connection
	vars := mux.Vars(r)
	boardID := vars["board"]
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
