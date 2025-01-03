package server

import (
	"log"
	"net/http"

	"github.com/ekaputra07/go-retro/internal/storage"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type WSServer struct {
	db              storage.Storage
	boards          map[*board]bool
	registerBoard   chan *board
	unregisterBoard chan *board
}

func (ws *WSServer) Start(stop chan struct{}) {
	log.Println("websocket server started")
	for {
		select {
		case b := <-ws.registerBoard:
			ws.boards[b] = true
			log.Printf("board=%s registered", b.ID)
		case b := <-ws.unregisterBoard:
			delete(ws.boards, b)
			log.Printf("board=%s unregistered", b.ID)
		case <-stop:
			log.Println("websocket server stopped")
			return
		}
	}
}

func (ws *WSServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	boardID := vars["board"]

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// get or start board process
	board := ws.getOrStartBoard(uuid.MustParse(boardID))

	// start user process
	user := newUser(conn)
	defer func() { board.leave <- user }()

	board.join <- user
	user.start()
}

func (ws *WSServer) getOrStartBoard(id uuid.UUID) *board {
	// if board is running, return it
	for b := range ws.boards {
		if b.ID == id {
			log.Printf("board=%s still running\n", b.ID)
			return b
		}
	}

	board, _ := getOrCreateBoard(id, ws)
	ws.registerBoard <- board
	go board.start()
	return board
}

func NewWSServer(db storage.Storage) *WSServer {
	return &WSServer{
		db:              db,
		boards:          make(map[*board]bool),
		registerBoard:   make(chan *board),
		unregisterBoard: make(chan *board),
	}
}
