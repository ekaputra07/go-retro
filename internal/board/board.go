package board

import (
	"errors"
	"fmt"
	"log"

	"github.com/ekaputra07/go-retro/internal/model"
	"github.com/ekaputra07/go-retro/internal/storage"
	"github.com/google/uuid"
)

// initial columns assigned when the board created
var defaultColumns = []string{"Good", "Bad", "Questions", "Emoji"}

// Board represents a single board instance that can be joined by clients
type Board struct {
	*model.Board
	manager *BoardManager
	db      storage.Storage
	clients map[*Client]bool
	timer   *timer

	// client joined and leaved
	join  chan *Client
	leave chan *Client

	// message to broadcast
	message chan *model.Message

	// to stop the board
	stop chan struct{}
}

// Add adds client to the board
func (b *Board) Add(client *Client) {
	b.join <- client
}

// Remove removes client from board
func (b *Board) Remove(client *Client) {
	b.leave <- client
}

func (b *Board) start() {
	// start the timer
	go b.timer.run()

	log.Printf("board=%s started", b.ID)
	for {
		select {
		case client := <-b.join:
			b.addClient(client)
			b.broadcastStatus()
			b.broadcastTimer()

		case client := <-b.leave:
			b.removeClient(client)
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
			b.manager.unregisterBoard <- b
			log.Printf("board=%s stopped", b.ID)
			return
		}
	}
}

func (b *Board) addClient(client *Client) {
	log.Printf("client=%s joined board=%s\n", client.ID, b.ID)
	b.clients[client] = true
}

func (b *Board) removeClient(client *Client) {
	if _, ok := b.clients[client]; ok {
		log.Printf("client=%s leaving board=%s\n", client.ID, b.ID)

		delete(b.clients, client)

		// if no joined clients, stop board
		if len(b.clients) == 0 {
			close(b.stop)
		}
	}
}

// update the board and broadcast its status if desired by returning `true` in bool output
func (b *Board) update(msg *model.Message) (bool, error) {
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

func (b *Board) broadcastStatus() error {
	// list clients
	var clients []*Client
	for c := range b.clients {
		clients = append(clients, c)
	}

	// list columns
	columns, err := b.db.ListColumn(b.ID)
	if err != nil {
		return fmt.Errorf("broadcastStatus failed while fetching columns: %s", err)
	}

	// list cards
	cards, err := b.db.ListCard(b.ID)
	if err != nil {
		return fmt.Errorf("broadcastStatus failed while fetching cards: %s", err)
	}
	msg := &model.Message{
		Type: model.MessageTypeBoardStatus,
		Data: map[string]any{
			"id":      b.ID,
			"clients": clients,
			"columns": columns,
			"cards":   cards,
		},
	}
	for u := range b.clients {
		u.message <- msg
	}
	return nil
}

func (b *Board) broadcastTimer() {
	msg := &model.Message{
		Type: model.MessageTypeTimerState,
		Data: b.timer,
	}
	for u := range b.clients {
		u.message <- msg
	}
}

func getOrCreateBoard(id uuid.UUID, manager *BoardManager) (*Board, error) {
	b, err := manager.db.GetBoard(id)

	// if board not in DB, create new one
	if err != nil {
		// these store operation should run in transaction
		b, err = manager.db.CreateBoard(id)
		if err != nil {
			return nil, err
		}
		for _, c := range defaultColumns {
			_, err := manager.db.CreateColumn(c, b.ID)
			if err != nil {
				return nil, err
			}
		}
	}

	return &Board{
		Board:   b,
		manager: manager,
		db:      manager.db,
		clients: make(map[*Client]bool),
		join:    make(chan *Client),
		leave:   make(chan *Client),
		message: make(chan *model.Message),
		stop:    make(chan struct{}),
		timer:   newTimer(),
	}, nil
}

func (b *Board) handleTimer(msg *model.Message) error {
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
