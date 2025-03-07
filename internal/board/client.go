package board

import (
	"log"
	"math/rand"
	"time"

	"github.com/ekaputra07/go-retro/internal/storage"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Client represents websocket connection between client (browser) that join a board
type Client struct {
	ID       uuid.UUID     `json:"id"`
	User     *storage.User `json:"user"`
	JoinedAt int64         `json:"joined_at"`
	AvatarID int           `json:"avatar_id"`

	board   *Board
	conn    *websocket.Conn
	message chan *message
}

// read reads message from socket
func (c *Client) read() {
	defer c.conn.Close()
	for {
		var msg message
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

// Start starts the client read (goroutine) and write (blocking) process
func (c *Client) Start() {
	log.Printf("client=%s started\n", c.ID)
	go c.write()

	// read message from connection
	c.read()
}

// Stop stops the client by closing the message channel
func (c *Client) Stop() {
	log.Printf("client=%s stopped\n", c.ID)
	close(c.message)
}

// NewClient creates a new client instance
func NewClient(conn *websocket.Conn, user *storage.User, board *Board) *Client {
	return &Client{
		ID:       uuid.New(),
		User:     user,
		JoinedAt: time.Now().Unix(),
		AvatarID: rand.Intn(11) + 1, // TODO: make this configurable
		board:    board,
		conn:     conn,
		message:  make(chan *message),
	}
}
