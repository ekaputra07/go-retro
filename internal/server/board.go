package server

import (
	"log"
	"time"

	"github.com/ekaputra07/go-retro/internal/model"
	"github.com/google/uuid"
)

type board struct {
	ID        string    `json:"id"`
	Columns   []*column `json:"columns"`
	Cards     []*card   `json:"cards"`
	UserCount int       `json:"user_count"`
	CreatedAt int64     `json:"created_at"`

	users   map[*user]bool
	join    chan *user
	leave   chan *user
	message chan *model.Message
}

func (b *board) Start() {
	log.Printf("board id=%s started", b.ID)
	for {
		select {
		case user := <-b.join:
			b.addUser(user)
		case user := <-b.leave:
			b.removeUser(user)
		case msg := <-b.message:
			b.updateAndBroadcast(msg)
		}
	}
}

func (b *board) addUser(user *user) {
	log.Printf("user id=%s joined board=%s\n", user.ID, b.ID)
	b.users[user] = true
	b.UserCount += 1
	user.board = b

	// notify users
	b.broadcastStatus()
}

func (b *board) removeUser(user *user) {
	if _, ok := b.users[user]; ok {
		log.Printf("user id=%s leaving board=%s\n", user.ID, b.ID)
		delete(b.users, user)
		b.UserCount -= 1

		// notify users
		b.broadcastStatus()
	}
}

func (b *board) updateAndBroadcast(msg *model.Message) {
	switch msg.Type {
	case model.MessageTypeCardNew:
		b.createCard(msg)
	case model.MessageTypeCardDelete:
		b.deleteCard(msg)
	case model.MessageTypeCardUpdate:
		b.updateCard(msg)
	}
}

func (b *board) broadcastStatus() {
	msg := &model.Message{
		Type: model.MessageTypeBoardStatus,
		Data: b,
	}
	for u := range b.users {
		u.message <- msg
	}
}

func newColumn(name string) *column {
	return &column{ID: uuid.New(), Name: name}
}

func newBoard(id string) *board {
	columns := []*column{}
	for _, c := range defaultColumns {
		columns = append(columns, &column{ID: uuid.New(), Name: c})
	}
	return &board{
		ID:        id,
		Columns:   columns,
		UserCount: 0,
		CreatedAt: time.Now().Unix(),
		users:     make(map[*user]bool),
		join:      make(chan *user),
		leave:     make(chan *user),
		message:   make(chan *model.Message),
	}
}
