package board

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"slices"

	"github.com/ekaputra07/go-retro/internal/models"
	"github.com/ekaputra07/go-retro/internal/natsutil"
	"github.com/ekaputra07/go-retro/internal/store"
	"github.com/ekaputra07/go-retro/internal/store/memstore"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

func boardStatusTopic(boardID uuid.UUID) string {
	return fmt.Sprintf("board.%s.status", boardID)
}

func clientJoinTopic(boardID uuid.UUID) string {
	return fmt.Sprintf("board.%s.client-joined", boardID)
}

func clientLeaveTopic(boardID uuid.UUID) string {
	return fmt.Sprintf("board.%s.client-leave", boardID)
}

func inboundMessageTopic(boardID uuid.UUID) string {
	return fmt.Sprintf("board.%s.msg-in", boardID)
}

func broadcastMessageTopic(boardID uuid.UUID) string {
	return fmt.Sprintf("board.%s.msg-all", boardID)
}

// Board represents a single board instance that can be joined by clients
type Board struct {
	*models.Board

	manager *BoardManager
	logger  *slog.Logger
	store   *store.BoardStore
	nats    *natsutil.NATS
	timer   *timer
	clients map[uuid.UUID]Client
	stop    chan bool

	// subscriptions channels
	joinCh    chan *nats.Msg
	leaveCh   chan *nats.Msg
	messageCh chan *nats.Msg
	statusCh  chan *nats.Msg
}

// Start starts the board and board's timer
func (b *Board) Start() {
	ctx, cancel := context.WithCancel(context.Background())

	// start the timer
	go b.timer.run(ctx)

	// listen to board events
	go b.listen(ctx, cancel)

	b.logger.Info("board started", "id", b.ID)
}

// Stop stops the board
func (b *Board) Stop() {
	b.stop <- true
}

func (b *Board) listen(ctx context.Context, cancel context.CancelFunc) {
	// create subscriptions
	joinSub, _ := b.nats.Conn.ChanSubscribe(clientJoinTopic(b.ID), b.joinCh)
	leaveSub, _ := b.nats.Conn.ChanSubscribe(clientLeaveTopic(b.ID), b.leaveCh)
	messageSub, _ := b.nats.Conn.ChanSubscribe(inboundMessageTopic(b.ID), b.messageCh)
	statusSub, _ := b.nats.Conn.ChanSubscribe(boardStatusTopic(b.ID), b.statusCh)

	defer func() {
		joinSub.Unsubscribe()
		leaveSub.Unsubscribe()
		messageSub.Unsubscribe()
		statusSub.Unsubscribe()

		close(b.joinCh)
		close(b.leaveCh)
		close(b.messageCh)
		close(b.statusCh)
	}()

	for {
		select {
		case msg := <-b.statusCh:
			msg.Respond(nil)

		case msg := <-b.joinCh:
			var c Client
			if err := json.Unmarshal(msg.Data, &c); err == nil {
				user := *c.User
				b.addClient(c)
				b.broadcast(b.usersStateMessage(user))
				b.broadcast(b.boardStateMessage(ctx, user))

				statuses := []timerStatus{timerStatusRunning, timerStatusPaused}
				if slices.Contains(statuses, b.timer.Status) {
					b.broadcast(b.timerStateMessage())
				}
			}

		case msg := <-b.leaveCh:
			var c Client
			if err := json.Unmarshal(msg.Data, &c); err == nil {
				b.removeClient(c)
				b.broadcast(b.usersStateMessage(*c.User))
			}

		case msg := <-b.messageCh:
			var m message
			if err := json.Unmarshal(msg.Data, &m); err == nil {
				broadcast, err := b.update(ctx, m)
				if err != nil {
					b.logger.Error("updating board failed", "board", b.ID, "err", err.Error())
					continue
				}

				if broadcast {
					b.broadcast(b.boardStateMessage(ctx, m.User))
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

func (b *Board) addClient(client Client) {
	b.logger.Info("client join board", "board", b.ID, "client", client.ID)
	b.clients[client.ID] = client
}

func (b *Board) removeClient(client Client) {
	if _, ok := b.clients[client.ID]; ok {
		b.logger.Info("client leave board", "board", b.ID, "client", client.ID)
		delete(b.clients, client.ID)
	}
}

// update the board and broadcast its status if desired
// (bool, error) --> (broadcast?, error)
func (b *Board) update(ctx context.Context, msg message) (bool, error) {
	switch msg.Type {
	case messageTypeColumnNew:
		return true, b.createColumn(ctx, msg)
	case messageTypeColumnDelete:
		return true, b.deleteColumn(ctx, msg)
	case messageTypeColumnUpdate:
		return true, b.updateColumn(ctx, msg)
	case messageTypeCardNew:
		return true, b.createCard(ctx, msg)
	case messageTypeCardDelete:
		return true, b.deleteCard(ctx, msg)
	case messageTypeCardUpdate:
		return true, b.updateCard(ctx, msg)
	case messageTypeCardVote:
		return true, b.voteCard(ctx, msg)
	case messageTypeTimerCmd:
		return false, b.handleTimerCommand(msg)
	}
	return false, nil
}

// usersStateMessage builds and returns the users state message
func (b *Board) usersStateMessage(user models.User) message {
	// clients map to slice
	clients := []Client{}
	for _, c := range b.clients {
		clients = append(clients, c)
	}

	return message{
		Type: messageTypeBoardUsers,
		Data: clients,
		User: user,
	}
}

// boardStateMessage builds and returns the board status message
func (b *Board) boardStateMessage(ctx context.Context, user models.User) message {
	// list columns
	columns, err := b.store.Columns.List(ctx)
	if err != nil {
		b.logger.Error("failed fetching columns", "board", b.ID, "err", err.Error())
	}

	// list cards
	cards, err := b.store.Cards.List(ctx)
	if err != nil {
		b.logger.Error("failed fetching cards", "board", b.ID, "err", err.Error())
	}

	return message{
		Type: messageTypeBoardStatus,
		Data: map[string]any{
			"id":      b.ID,
			"columns": columns,
			"cards":   cards,
		},
		User: user,
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

// newBoard creates board instance (using in-memory store)
func newBoard(_ context.Context, manager *BoardManager, board *models.Board) (*Board, error) {
	// create new storage to store board's data
	// store, err := natstore.NewBoardStore(ctx, manager.nats, fmt.Sprintf("goretro-board-%s", board.ID))
	// if err != nil {
	// 	return nil, err
	// }
	store := memstore.NewBoardStore()

	return &Board{
		Board:     board,
		manager:   manager,
		logger:    manager.logger,
		store:     store,
		nats:      manager.nats,
		clients:   make(map[uuid.UUID]Client),
		timer:     newTimer(manager.logger),
		stop:      make(chan bool),
		joinCh:    make(chan *nats.Msg, 1),
		leaveCh:   make(chan *nats.Msg, 1),
		messageCh: make(chan *nats.Msg, 1),
		statusCh:  make(chan *nats.Msg, 1),
	}, nil
}
