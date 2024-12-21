package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type WSServer struct {
	boards map[*board]bool
}

func (ws *WSServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	boardID := vars["board"]

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	board := ws.getOrStartBoard(boardID)
	user := newUser(conn)
	defer func() { board.leave <- user }()

	board.join <- user
	user.start()
}

func (ws *WSServer) getOrStartBoard(id string) *board {
	var b *board
	for eb := range ws.boards {
		if eb.ID == id {
			log.Printf("existing board=%s found\n", eb.ID)
			return eb
		}
	}

	// not running, create and run new one
	b = newBoard(id)
	log.Printf("new board=%s created\n", b.ID)
	go b.Start()
	ws.boards[b] = true
	return b
}

func NewWSServer() *WSServer {
	return &WSServer{
		boards: make(map[*board]bool),
	}
}
