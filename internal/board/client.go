package board

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/ekaputra07/go-retro/internal/models"
	"github.com/ekaputra07/go-retro/internal/natsutil"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = 10 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// Client represents websocket connection between client (browser) that join a board
type Client struct {
	ID       uuid.UUID    `json:"id"`
	BoardID  uuid.UUID    `json:"board_id"`
	User     *models.User `json:"user"`
	JoinedAt int64        `json:"joined_at"`

	logger    *slog.Logger
	conn      *websocket.Conn
	nats      *natsutil.NATS
	messageCh chan *nats.Msg
}

func (c *Client) publish(topic string, msg any) {
	go func() {
		data, err := json.Marshal(msg)
		if err != nil {
			c.logger.Error(fmt.Sprintf("failed marshaling message: %s", err.Error()))
		}
		if err = c.nats.Conn.Publish(topic, data); err != nil {
			c.logger.Error(fmt.Sprintf("failed publishing message: %s", err.Error()))
		}
	}()
}

// read reads message from socket
func (c *Client) read() {
	defer c.conn.Close()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(appData string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		var msg message
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			c.logger.Error("client error reading <--", "id", c.ID, "err", err.Error())
			break
		}

		msg.User = *c.User

		// handle `me` message
		if msg.Type == messageTypeMe {
			data, _ := json.Marshal(msg)
			c.messageCh <- &nats.Msg{Data: data}
		} else {
			c.publish(inboundMessageTopic(c.BoardID), msg)
		}
	}
}

// write writes message to the socket
func (c *Client) write(ctx context.Context) {
	// subscribe for messages
	messageSub, _ := c.nats.Conn.ChanSubscribe(broadcastMessageTopic(c.BoardID), c.messageCh)

	// watch for columns and cards changes
	kv, _ := c.nats.JS.KeyValue(ctx, "goretro")
	w, _ := kv.WatchFiltered(ctx, []string{
		fmt.Sprintf("boards.%s.columns.*", c.BoardID),
		fmt.Sprintf("boards.%s.cards.*", c.BoardID),
	})

	// pinger
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		messageSub.Unsubscribe()
		w.Stop()
		close(c.messageCh)
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case kve := <-w.Updates():
			if kve != nil {
				s, err := newStreamFromKVE(kve)
				if s != nil && err == nil {
					c.conn.SetWriteDeadline(time.Now().Add(writeWait))
					if err := c.conn.WriteJSON(s); err != nil {
						c.logger.Error("client message error -->", "id", c.ID, "err", err.Error())
						return
					}
				}
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.logger.Error("client ping error -->", "id", c.ID, "err", err.Error())
				return
			}
		case msg := <-c.messageCh:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.TextMessage, msg.Data); err != nil {
				c.logger.Error("client message error -->", "id", c.ID, "err", err.Error())
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

// Start starts the client read (goroutine) and write (blocking) process
func (c *Client) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// start writer
	go c.write(ctx)

	// client is now ready, notify board
	c.publish(clientJoinTopic(c.BoardID), c)
	defer func() {
		c.publish(clientLeaveTopic(c.BoardID), c)
	}()
	c.logger.Info("client started", "id", c.ID)

	// read message from connection
	c.read()
}

// NewClient creates a new client instance
func NewClient(conn *websocket.Conn,
	user *models.User,
	logger *slog.Logger,
	nats_ *natsutil.NATS,
	boardID uuid.UUID,
) *Client {
	return &Client{
		ID:        uuid.New(),
		User:      user,
		JoinedAt:  time.Now().Unix(),
		BoardID:   boardID,
		logger:    logger,
		conn:      conn,
		nats:      nats_,
		messageCh: make(chan *nats.Msg, 256),
	}
}
