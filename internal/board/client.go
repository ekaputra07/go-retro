package board

import (
	"log"
	"math/rand/v2"
	"time"

	"github.com/ekaputra07/go-retro/internal/model"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Client represents websocket connection between client (browser) that join a board
type Client struct {
	ID       uuid.UUID   `json:"id"`
	User     *model.User `json:"user"`
	JoinedAt int64       `json:"joined_at"`
	AvatarID int         `json:"avatar_id"`
	board    *Board
	conn     *websocket.Conn
	message  chan *Message
}

// read reads message from socket
func (c *Client) read() {
	defer c.conn.Close()
	for {
		var msg Message
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("client=%s error reading <-- %s", c.ID, err)
			break
		}
		log.Printf("client=%s <-- %s", c.ID, msg.Type)

		msg.client = c
		c.board.message <- &msg
	}
}

// write writes message to the socket
func (c *Client) write() {
	defer c.conn.Close()
	for msg := range c.message {
		if c == msg.client {
			continue
		}

		log.Printf("client=%s --> %v", c.ID, msg.Type)
		if err := c.conn.WriteJSON(msg); err != nil {
			log.Printf("client=%s error writing --> %v", c.ID, err)
			break
		}
	}
}

func (c *Client) Start() {
	log.Printf("client=%s started\n", c.ID)
	go c.write()

	// read message from connection
	c.read()
}

func (c *Client) Stop() {
	log.Printf("client=%s stopped\n", c.ID)
	close(c.message)
}

func NewClient(conn *websocket.Conn, user *model.User, board *Board) *Client {
	return &Client{
		ID:       uuid.New(),
		User:     user,
		JoinedAt: time.Now().Unix(),
		AvatarID: rand.IntN(11) + 1,
		board:    board,
		conn:     conn,
		message:  make(chan *Message),
	}
}
