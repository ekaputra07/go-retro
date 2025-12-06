package board

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/ekaputra07/go-retro/internal/models"
	"github.com/ekaputra07/go-retro/internal/natsutil"
	"github.com/ekaputra07/go-retro/internal/store"
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
	*models.Client

	logger     *slog.Logger
	conn       *websocket.Conn
	nats       *natsutil.NATS
	store      *store.Store
	msgHandler *messageHandler
	messageCh  chan *nats.Msg
}

// publish publish message to subscribers via nats
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

// checkTimerStateMessage check for latest state of active timer.
// only return message when its in 'running' or 'paused' state.
func (c *Client) checkTimerStateMessage() *message {
	msg, err := queryTimerStatus(c.nats.Conn, c.BoardID)
	if err != nil {
		c.logger.Error("error requesting timer status message", "err", err.Error())
		return nil
	}
	var m message
	if err = json.Unmarshal(msg.Data, &m); err != nil {
		c.logger.Error("error decoding timer status message", "err", err.Error())
		return nil
	}
	var status string
	err = m.stringVar(&status, "status")
	if err != nil {
		c.logger.Error("error reading status", "err", err.Error())
		return nil
	}
	// ignore timer state when its stopped or done.
	if slices.Contains([]string{string(timerStatusRunning), string(timerStatusPaused)}, status) {
		return &m
	}
	return nil
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
		msg.BoardID = c.BoardID

		switch msg.Type {
		case messageTypeMe:
			msgs := []message{msg}

			// during ME inqury, check timer state and includes in messages if any.
			// returning timer state to user is necessary so that new joined user could
			// see the timer UI even when it's currently paused since paused timer don't emit events.
			if timerStateMsg := c.checkTimerStateMessage(); timerStateMsg != nil {
				msgs = append(msgs, *timerStateMsg)
			}

			ml := newMessageList(c.BoardID, msgs...)
			data, err := ml.encode()
			if err != nil {
				c.logger.Error("failed to encode messageList during ME", "err", err.Error())
			}
			c.messageCh <- &nats.Msg{Data: data}
		case messageTypeTimerCmd:
			c.publish(timerCmdTopic(c.BoardID), msg)
		default:
			c.msgHandler.handle(context.Background(), msg)
		}
	}
}

// write writes message to the socket
func (c *Client) write(ctx context.Context) {
	// subscribe for messages
	messageSub, err := c.nats.Conn.ChanSubscribe(broadcastMessageTopic(c.BoardID), c.messageCh)
	if err != nil {
		c.logger.Error("client subscribe error -->", "id", c.ID, "err", err.Error())
		return
	}

	// watch for clients, columns and cards changes
	kv, err := c.nats.JS.KeyValue(ctx, "goretro")
	if err != nil {
		c.logger.Error("client kv error -->", "id", c.ID, "err", err.Error())
		return
	}
	w, err := kv.WatchFiltered(ctx, []string{
		fmt.Sprintf("boards.%s.clients.*", c.BoardID),
		fmt.Sprintf("boards.%s.columns.*", c.BoardID),
		fmt.Sprintf("boards.%s.cards.*", c.BoardID),
	})
	if err != nil {
		c.logger.Error("client watch error -->", "id", c.ID, "err", err.Error())
		return
	}

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
				s, err := newStream(kve.Key(), kve.Operation(), kve.Value())
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

	defer func() {
		// delete client on leave
		err := c.store.Clients.Delete(ctx, c.BoardID, c.ID)
		if err != nil {
			c.logger.Error("error deleting client record", "board", c.BoardID, "id", c.ID)
		}
	}()
	c.logger.Info("client started", "id", c.ID)

	// read message from connection
	c.read()
}

// NewClient creates a new client instance
func NewClient(
	ctx context.Context,
	conn *websocket.Conn,
	user *models.User,
	logger *slog.Logger,
	store *store.Store,
	nats_ *natsutil.NATS,
	boardID uuid.UUID,
) (*Client, error) {
	// create client record
	model := models.NewClient(user, boardID)
	err := store.Clients.Create(ctx, model)
	if err != nil {
		return nil, err

	}
	// create client process instance
	return &Client{
		Client:     &model,
		logger:     logger,
		conn:       conn,
		nats:       nats_,
		store:      store,
		msgHandler: newMessageHandler(store),
		messageCh:  make(chan *nats.Msg, 256),
	}, nil
}
