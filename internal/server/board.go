package server

import (
	"errors"
	"fmt"
	"log"

	"github.com/ekaputra07/go-retro/internal/model"
	"github.com/ekaputra07/go-retro/internal/storage"
	"github.com/google/uuid"
)

type board struct {
	*model.Board
	UserCount int `json:"user_count"`
	ws        *WSServer
	db        storage.Storage
	users     map[*user]bool
	timer     *timer

	// user joined and leaved
	join  chan *user
	leave chan *user

	// message to broadcast
	message chan *model.Message

	// to stop the board
	stop chan struct{}
}

func (b *board) start() {
	// start the timer
	go b.timer.run()

	log.Printf("board=%s started", b.ID)
	for {
		select {
		case user := <-b.join:
			b.addUser(user)
			b.broadcastStatus()
			b.broadcastTimer()

		case user := <-b.leave:
			b.removeUser(user)
			b.broadcastStatus()

		case msg := <-b.message:
			broadcast, err := b.update(msg)
			if err != nil {
				log.Printf("updating board=%s failed: %s", b.ID, err)
				continue
			}
			if broadcast {
				if err := b.broadcastStatus(); err != nil {
					log.Printf("broadcasting board=%s status failed: %s", b.ID, err)
				}
			}

		case <-b.timer.state:
			b.broadcastTimer()

		case <-b.stop:
			// cleanup timer when board stopped
			if b.timer != nil {
				b.timer.cmd <- timerCmd{cmd: "destroy"}
				b.timer = nil
			}

			// stop and unregister from ws server
			b.ws.unregisterBoard <- b
			log.Printf("board=%s stopped", b.ID)
			return
		}
	}
}

func (b *board) addUser(user *user) {
	log.Printf("user=%s joined board=%s\n", user.ID, b.ID)
	b.users[user] = true
	b.UserCount++
	user.board = b
}

func (b *board) removeUser(user *user) {
	if _, ok := b.users[user]; ok {
		log.Printf("user=%s leaving board=%s\n", user.ID, b.ID)

		user.stop()
		delete(b.users, user)
		b.UserCount--

		// if no joined users, stop board
		if b.UserCount == 0 {
			close(b.stop)
		}
	}
}

// update the board and broadcast its status if desired by returning `true` in bool output
func (b *board) update(msg *model.Message) (bool, error) {
	switch msg.Type {
	case model.MessageTypeColumnNew:
		return true, b.createColumn(msg)
	case model.MessageTypeColumnDelete:
		return true, b.deleteColumn(msg)
	case model.MessageTypeColumnUpdate:
		return true, b.updateColumn(msg)
	case model.MessageTypeCardNew:
		return true, b.createCard(msg)
	case model.MessageTypeCardDelete:
		return true, b.deleteCard(msg)
	case model.MessageTypeCardUpdate:
		return true, b.updateCard(msg)
	case model.MessageTypeCardVote:
		return true, b.voteCard(msg)
	case model.MessageTypeTimerCmd:
		return false, b.handleTimer(msg)
	}
	return false, nil
}

func (b *board) broadcastStatus() error {
	columns, err := b.db.ListColumn(b.ID)
	if err != nil {
		return fmt.Errorf("broadcastStatus failed while fetching columns: %s", err)
	}
	cards, err := b.db.ListCard(b.ID)
	if err != nil {
		return fmt.Errorf("broadcastStatus failed while fetching cards: %s", err)
	}
	msg := &model.Message{
		Type: model.MessageTypeBoardStatus,
		Data: map[string]interface{}{
			"id":         b.ID,
			"user_count": b.UserCount,
			"columns":    columns,
			"cards":      cards,
		},
	}
	for u := range b.users {
		u.message <- msg
	}
	return nil
}

func (b *board) broadcastTimer() {
	msg := &model.Message{
		Type: model.MessageTypeTimerState,
		Data: b.timer,
	}
	for u := range b.users {
		u.message <- msg
	}
}

func getOrCreateBoard(id uuid.UUID, ws *WSServer) (*board, error) {
	b, err := ws.db.GetBoard(id)

	// if board not in DB, create new one
	if err != nil {
		// these store operation should run in transaction
		b, err = ws.db.CreateBoard(id)
		if err != nil {
			return nil, err
		}
		for _, c := range defaultColumns {
			_, err := ws.db.CreateColumn(c, b.ID)
			if err != nil {
				return nil, err
			}
		}
	}

	return &board{
		Board:     b,
		UserCount: 0,
		ws:        ws,
		db:        ws.db,
		users:     make(map[*user]bool),
		join:      make(chan *user),
		leave:     make(chan *user),
		message:   make(chan *model.Message),
		stop:      make(chan struct{}),
		timer:     newTimer(),
	}, nil
}

func (b *board) handleTimer(msg *model.Message) error {
	data := msg.Data.(map[string]any)
	cmdAny, ok := data["cmd"]
	if !ok {
		return errors.New("handleTimer payload missing `cmd` field")
	}
	cmd := timerCmd{cmd: cmdAny.(string)}

	value, ok := data["value"]
	if ok {
		cmd.value = value.(string)
	}

	b.timer.cmd <- cmd
	return nil
}
