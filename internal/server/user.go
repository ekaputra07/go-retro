package server

import (
	"log"

	"github.com/ekaputra07/go-retro/internal/model"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type user struct {
	ID      uuid.UUID `json:"id"`
	conn    *websocket.Conn
	board   *board
	message chan *model.Message
}

func (u *user) read() {
	defer u.conn.Close()
	for {
		var msg model.Message
		err := u.conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("user=%s error reading <-- %s", u.ID, err)
			break
		}
		log.Printf("user=%s <-- %s", u.ID, msg.Type)
		u.board.message <- &msg
	}
}

func (u *user) write() {
	defer u.conn.Close()
	for msg := range u.message {
		log.Printf("user=%s --> %v", u.ID, msg.Type)
		if err := u.conn.WriteJSON(msg); err != nil {
			log.Printf("user=%s error writing --> %v", u.ID, err)
			break
		}
	}
}

func (u *user) start() {
	log.Printf("user=%s started\n", u.ID)
	go u.write()

	// read message from client
	u.read()
}

func (u *user) stop() {
	log.Printf("user=%s stopped\n", u.ID)
	close(u.message)
}

func newUser(conn *websocket.Conn) *user {
	return &user{
		ID:      uuid.New(),
		conn:    conn,
		message: make(chan *model.Message),
	}
}
