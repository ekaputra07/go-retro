package board

import (
	"log/slog"
	"time"

	"github.com/ekaputra07/go-retro/internal/models"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Client represents websocket connection between client (browser) that join a board
type Client struct {
	ID       uuid.UUID    `json:"id"`
	User     *models.User `json:"user"`
	JoinedAt int64        `json:"joined_at"`

	board   *Board
	logger  *slog.Logger
	conn    *websocket.Conn
	message chan message
}

// read reads message from socket
func (c *Client) read() {
	defer c.conn.Close()
	for {
		var msg message
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			c.logger.Error("client error reading <--", "id", c.ID, "err", err.Error())
			break
		}
		c.logger.Info("client <--", "id", c.ID, "type", msg.Type)

		msg.fromClient = c
		c.board.message <- msg
	}
}

// write writes message to the socket
func (c *Client) write() {
	defer c.conn.Close()

	ticker := time.NewTicker(5 * time.Second)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case <-ticker.C:
			if err := c.conn.WriteJSON(message{Type: messageTypePing}); err != nil {
				c.logger.Error("client ping error -->", "id", c.ID, "err", err.Error())
				return
			}
		case msg := <-c.message:
			c.logger.Info("client -->", "id", c.ID, "type", msg.Type)
			if err := c.conn.WriteJSON(msg); err != nil {
				c.logger.Error("client message error -->", "id", c.ID, "err", err.Error())
				return
			}
		}
	}
}

// Start starts the client read (goroutine) and write (blocking) process
func (c *Client) Start() {
	c.logger.Info("client started", "id", c.ID)
	go c.write()

	// read message from connection
	c.read()
}

// Stop stops the client by closing the message channel
func (c *Client) Stop() {
	c.logger.Info("client stopped", "id", c.ID)
	close(c.message)
}

// NewClient creates a new client instance
func NewClient(conn *websocket.Conn, user *models.User, board *Board) *Client {
	return &Client{
		ID:       uuid.New(),
		User:     user,
		JoinedAt: time.Now().Unix(),
		board:    board,
		logger:   board.logger,
		conn:     conn,
		message:  make(chan message),
	}
}
