package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ekaputra07/go-retro/internal/board"
	"github.com/ekaputra07/go-retro/internal/storage"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
)

func init() {
	gob.Register(uuid.UUID{})
}

const SESSION_NAME = "goretro_session"

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type app struct {
	db      storage.Storage
	manager *board.BoardManager
	router  *mux.Router
	session *sessions.CookieStore
}

func (a *app) start() context.CancelFunc {
	// session store
	secret := os.Getenv("GORETRO_SESSION_SECRET")
	secure := os.Getenv("GORETRO_SESSION_SECURE") != "false" // secure by default
	if secret == "" {
		log.Fatalln("GORETRO_SESSION_SECRET not set!")
	}
	a.session = sessions.NewCookieStore([]byte(secret))
	a.session.Options = &sessions.Options{Secure: secure}

	// start WS server
	a.manager = board.NewBoardManager(a.db)
	ctx, cancel := context.WithCancel(context.Background())
	go a.manager.Start(ctx)

	// setup routes
	a.router = mux.NewRouter()
	a.router.Use(a.loggingMiddleware)
	a.router.HandleFunc("/health", a.health)
	a.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web"))))
	a.router.HandleFunc("/b/{board}/ws", a.websocket)
	a.router.HandleFunc("/b/{board}", a.board)
	a.router.HandleFunc("/", a.generateBoard)

	return cancel
}

func (a *app) loggingMiddleware(next http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, next)
}

func (a *app) health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session.Values["user_id"] = u.ID
		if err := session.Save(r, w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("new user created with id=%s", u.ID)
	}

	http.ServeFile(w, r, "web/index.html")
}

func (a *app) websocket(w http.ResponseWriter, r *http.Request) {
	// validate session (make sure user is present) before upgrading the connection
	// TODO: Move this check to a middleware
	session, _ := a.session.Get(r, SESSION_NAME)
	userID := session.Values["user_id"].(uuid.UUID)

	user, err := a.db.GetUser(userID)
	if err != nil {
		http.Error(w, "unauthorized access", http.StatusUnauthorized)
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
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
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
