package board

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/ekaputra07/go-retro/internal/models"
	"github.com/ekaputra07/go-retro/internal/natsutil"
	"github.com/ekaputra07/go-retro/internal/store"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

func boardStatusTopic(boardID uuid.UUID) string {
	return fmt.Sprintf("boards.%s.status", boardID)
}

func clientJoinTopic(boardID uuid.UUID) string {
	return fmt.Sprintf("boards.%s.client-joined", boardID)
}

func clientLeaveTopic(boardID uuid.UUID) string {
	return fmt.Sprintf("boards.%s.client-leave", boardID)
}

func inboundMessageTopic(boardID uuid.UUID) string {
	return fmt.Sprintf("boards.%s.msg-in", boardID)
}

func broadcastMessageTopic(boardID uuid.UUID) string {
	return fmt.Sprintf("boards.%s.msg-out", boardID)
}

// Board represents a single board instance that can be joined by clients
type Board struct {
	*models.Board

	manager    *BoardManager
	logger     *slog.Logger
	store      *store.Store
	nats       *natsutil.NATS
	timer      *timer
	msgHandler *messageHandler
	clients    map[uuid.UUID]Client
	stop       chan bool

	// subscriptions channels
	messageCh chan *nats.Msg
	statusCh  chan *nats.Msg
}

// Start starts the board and board's timer
func (b *Board) Start() {
	ctx, cancel := context.WithCancel(context.Background())

	// start the timer
	go b.timer.run(ctx)

	// listen to board events
	go b.listen(cancel)

	b.logger.Info("board started", "id", b.ID)
}

// Stop stops the board
func (b *Board) Stop() {
	b.stop <- true
}

func (b *Board) listen(cancel context.CancelFunc) {
	// create subscriptions
	messageSub, _ := b.nats.Conn.ChanSubscribe(inboundMessageTopic(b.ID), b.messageCh)
	statusSub, _ := b.nats.Conn.ChanSubscribe(boardStatusTopic(b.ID), b.statusCh)

	defer func() {
		messageSub.Unsubscribe()
		statusSub.Unsubscribe()
		close(b.messageCh)
		close(b.statusCh)
	}()

	for {
		select {
		case msg := <-b.statusCh:
			msg.Respond(nil)

		case msg := <-b.messageCh:
			var m message
			if err := json.Unmarshal(msg.Data, &m); err == nil {
				// board no more handle CRUD message (moved to clients) so now only handle timer command here.
				if m.Type == messageTypeTimerCmd {
					if err := b.handleTimerCommand(m); err != nil {
						b.logger.Error("handleTimerCommand failed", "board", b.ID, "err", err.Error())
						continue
					}
				}
			}

		case t := <-b.timer.state:
			// broadcast timer state and notify timer state change
			b.broadcast(b.timerStateMessage())

			// if timer state changes was a result of user action: broadcast notification
			if t.lastCmdUser.ID != uuid.Nil {
				b.broadcast(b.notificationMessage(t.statusMessage, t.lastCmdUser))
			}

		case <-b.stop:
			// stop timer
			cancel()

			// stop and unregister from ws server
			b.manager.unregisterChan <- b
			b.logger.Info("board stopped", "board", b.ID)
			return
		}
	}
}

// timerStateMessage builds and returns the timer state message
func (b *Board) timerStateMessage() message {
	return message{
		Type: messageTypeTimerState,
		Data: b.timer,
	}
}

// notificationMessage builds and returns the notification message
func (b *Board) notificationMessage(msg string, user models.User) message {
	return message{
		Type: messageTypeBoardNotification,
		Data: msg,
		User: user,
	}
}

func (b *Board) broadcast(msg message) {
	data, err := json.Marshal(msg)
	if err != nil {
		b.logger.Error(fmt.Sprintf("failed marshaling message: %s", err.Error()))
	}
	if err = b.nats.Conn.Publish(broadcastMessageTopic(b.ID), data); err != nil {
		b.logger.Error(fmt.Sprintf("failed publishing message: %s", err.Error()))
	}
}

// handleTimerCommand handles timer command message
func (b *Board) handleTimerCommand(msg message) error {
	var cmd string
	var value string

	if err := msg.stringVar(&cmd, "cmd"); err != nil {
		return err
	}
	command := timerCmd{cmd: cmd, user: msg.User}

	if err := msg.stringVar(&value, "value"); err == nil {
		command.value = value
	}
	b.timer.cmd <- command
	return nil
}

// newBoard creates board instance
func newBoard(manager *BoardManager, board *models.Board) (*Board, error) {
	return &Board{
		Board:      board,
		manager:    manager,
		logger:     manager.logger,
		store:      manager.store,
		nats:       manager.nats,
		clients:    make(map[uuid.UUID]Client),
		timer:      newTimer(manager.logger),
		msgHandler: newMessageHandler(manager.store),
		stop:       make(chan bool),
		messageCh:  make(chan *nats.Msg, 256),
		statusCh:   make(chan *nats.Msg, 256),
	}, nil
}
