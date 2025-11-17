package board

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"

	"github.com/ekaputra07/go-retro/internal/models"
	"github.com/ekaputra07/go-retro/internal/nats"
	"github.com/ekaputra07/go-retro/internal/store"
	nstore "github.com/ekaputra07/go-retro/internal/store/nats"
)

// Board represents a single board instance that can be joined by clients
type Board struct {
	*models.Board

	manager *BoardManager
	logger  *slog.Logger
	store   *store.BoardStore
	clients map[*Client]bool
	timer   *timer

	// client joined and leaved
	join  chan *Client
	leave chan *Client

	// message to broadcast
	message chan message

	// to stop the board
	stop chan struct{}
}

// Add adds client to the board
func (b *Board) AddClient(client *Client) {
	b.join <- client
}

// Remove removes client from board
func (b *Board) RemoveClient(client *Client) {
	b.leave <- client
}

// Start starts the board and board's timer
func (b *Board) Start() {
	b.logger.Info("board started", "board", b.ID)

	// use context to stop timer from board process
	ctx, cancelTimer := context.WithCancel(context.Background())

	// start the timer
	go b.timer.run(ctx)

	// listen to board events
	go b.listen(cancelTimer)
}

func (b *Board) listen(cancelTimer context.CancelFunc) {
	ctx := context.Background()

	for {
		select {
		case client := <-b.join:
			b.addClient(client)

			us := b.usersStateMessage()
			bs := b.boardStateMessage(ctx)
			no := b.notificationMessage(fmt.Sprintf("%s joined", client.User.Name))

			// broadcast user state and notification to all except to newly joined client
			b.broadcast([]message{us, no}, client.User)

			// send joined user all messages except notification
			// also send timer state only when its running or paused
			msgs := []message{us, bs}
			if slices.Contains([]timerStatus{timerStatusRunning, timerStatusPaused}, b.timer.Status) {
				msgs = append(msgs, b.timerStateMessage())
			}
			b.send(msgs, client.User)

		case client := <-b.leave:
			b.removeClient(client)

			// broadcast board status and leave notification
			msgs := []message{
				b.usersStateMessage(),
				b.boardStateMessage(ctx),
			}
			b.broadcast(msgs, nil)

		case msg := <-b.message:
			broadcast, err := b.update(ctx, msg)
			if err != nil {
				b.logger.Error("updating board failed", "board", b.ID, "err", err.Error())
				continue
			}

			// broadcast board status if update is successful
			if broadcast {
				b.broadcast([]message{b.boardStateMessage(ctx)}, nil)
			}

		case t := <-b.timer.state:
			// broadcast timer state and notify timer state change
			ts := b.timerStateMessage()
			no := b.notificationMessage(t.statusMessage)

			// broadcast timer state to all
			b.broadcast([]message{ts}, nil)

			// if timer state changes was a result of user action (commandClient)
			// broadcast notification to all except lastCommandClient
			if t.lastCommandClient != nil {
				b.broadcast([]message{no}, t.lastCommandClient.User)
			}

		case <-b.stop:
			// cleanup timer when board stopped
			if b.timer != nil {
				cancelTimer()
			}

			// stop and unregister from ws server
			b.manager.unregisterChan <- b
			b.logger.Info("board stopped", "board", b.ID)
			return
		}
	}
}

func (b *Board) addClient(client *Client) {
	b.logger.Info("client join board", "board", b.ID, "client", client.ID)
	b.clients[client] = true
}

func (b *Board) removeClient(client *Client) {
	if _, ok := b.clients[client]; ok {
		b.logger.Info("client leave board", "board", b.ID, "client", client.ID)
		delete(b.clients, client)
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
func (b *Board) usersStateMessage() message {
	// clients map to slice
	var clients []*Client
	for c := range b.clients {
		clients = append(clients, c)
	}

	return message{
		Type: messageTypeBoardUsers,
		Data: clients,
	}
}

// boardStateMessage builds and returns the board status message
func (b *Board) boardStateMessage(ctx context.Context) message {
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
func (b *Board) notificationMessage(msg string) message {
	return message{
		Type: messageTypeBoardNotification,
		Data: msg,
	}
}

// send sends messages to specific user
func (b *Board) send(msgs []message, user *models.User) {
	for c := range b.clients {
		if user != nil && c.User.ID == user.ID {
			for _, m := range msgs {
				c.message <- m
			}
		}

	}
}

// broadcast sends messages to all users except excludeUser
func (b *Board) broadcast(msgs []message, excludeUser *models.User) {
	for c := range b.clients {
		if excludeUser != nil && c.User.ID == excludeUser.ID {
			continue
		}
		for _, m := range msgs {
			c.message <- m
		}
	}
}

// handleTimerCommand handles timer command message
func (b *Board) handleTimerCommand(msg message) error {
	data := msg.Data.(map[string]any)
	cmdAny, ok := data["cmd"]
	if !ok {
		return errors.New("handleTimerCommand payload missing `cmd` field")
	}
	cmd := timerCmd{cmd: cmdAny.(string), client: msg.fromClient}

	value, ok := data["value"]
	if ok {
		cmd.value = value.(string)
	}

	b.timer.cmd <- cmd
	return nil
}

// newBoard creates board instance (using in-memory store)
func newBoard(ctx context.Context, nats *nats.NATS, manager *BoardManager, board *models.Board) (*Board, error) {
	store, err := nstore.NewBoardStore(ctx, nats, fmt.Sprintf("goretro-board-%s", board.ID))
	if err != nil {
		return nil, err
	}
	return &Board{
		Board:   board,
		manager: manager,
		logger:  manager.logger,
		store:   store,
		clients: make(map[*Client]bool),
		join:    make(chan *Client),
		leave:   make(chan *Client),
		message: make(chan message),
		stop:    make(chan struct{}),
		timer:   newTimer(manager.logger),
	}, nil
}
