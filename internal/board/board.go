package board

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"slices"

	"github.com/ekaputra07/go-retro/internal/storage"
	"github.com/google/uuid"
)

// avatarsCount is the total number of avatars available to choose from.
// see: web/public/avatars
const avatarsCount = 12

// initial columns assigned when the board created
var defaultColumns = []string{"Good", "Bad", "Questions", "Emoji"}

// Board represents a single board instance that can be joined by clients
type Board struct {
	*storage.Board
	manager *BoardManager
	db      storage.Storage
	clients map[*Client]bool
	avatars map[int]bool
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
	log.Printf("board=%s started", b.ID)

	// start the timer
	go b.timer.run()

	// listen to board events
	go b.listen()
}

func (b *Board) listen() {
	for {
		select {
		case client := <-b.join:
			b.addClient(client)

			us := b.usersStateMessage()
			bs := b.boardStateMessage()
			no := b.notificationMessage(fmt.Sprintf("%s joined", client.User.Name))

			// broadcast user state and notification to all except to newly joined client
			b.broadcast([]message{us, no}, client)

			// send joined client all messages except notification
			// also send timer state only when its running or paused
			msgs := []message{us, bs}
			if slices.Contains([]timerStatus{timerStatusRunning, timerStatusPaused}, b.timer.Status) {
				msgs = append(msgs, b.timerStateMessage())
			}
			b.send(msgs, client)

		case client := <-b.leave:
			b.removeClient(client)

			// broadcast board status and leave notification
			msgs := []message{
				b.usersStateMessage(),
				b.boardStateMessage(),
			}
			b.broadcast(msgs, nil)

		case msg := <-b.message:
			broadcast, err := b.update(msg)
			if err != nil {
				log.Printf("updating board=%s failed: %s", b.ID, err)
				continue
			}

			// broadcast board status if update is successful
			if broadcast {
				b.broadcast([]message{b.boardStateMessage()}, nil)
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
				b.broadcast([]message{no}, t.lastCommandClient)
			}

		case <-b.stop:
			// cleanup timer when board stopped
			if b.timer != nil {
				b.timer.cmd <- timerCmd{cmd: "destroy"}
				b.timer = nil
			}

			// stop and unregister from ws server
			b.manager.unregisterChan <- b
			log.Printf("board=%s stopped", b.ID)
			return
		}
	}
}

func (b *Board) addClient(client *Client) {
	log.Printf("client=%s joined board=%s\n", client.ID, b.ID)
	avatarID := b.uniqueAvatarID()
	// TODO: avatar should be assigned to user not client
	client.AvatarID = avatarID
	b.avatars[avatarID] = true
	b.clients[client] = true
}

func (b *Board) removeClient(client *Client) {
	if _, ok := b.clients[client]; ok {
		log.Printf("client=%s leaving board=%s\n", client.ID, b.ID)

		delete(b.clients, client)
		delete(b.avatars, client.AvatarID)

		// if no joined clients, stop board
		if len(b.clients) == 0 {
			close(b.stop)
		}
	}
}

// uniqueAvatarID generates unique avatar ID for a client
// by checking if the ID is already used by another client recursively
func (b *Board) uniqueAvatarID() int {
	randID := rand.Intn(avatarsCount-1) + 1
	if _, ok := b.avatars[randID]; ok {
		return b.uniqueAvatarID()
	}
	return randID
}

// update the board and broadcast its status if desired
// (bool, error) --> (broadcast?, error)
func (b *Board) update(msg message) (bool, error) {
	switch msg.Type {
	case messageTypeColumnNew:
		return true, b.createColumn(msg)
	case messageTypeColumnDelete:
		return true, b.deleteColumn(msg)
	case messageTypeColumnUpdate:
		return true, b.updateColumn(msg)
	case messageTypeCardNew:
		return true, b.createCard(msg)
	case messageTypeCardDelete:
		return true, b.deleteCard(msg)
	case messageTypeCardUpdate:
		return true, b.updateCard(msg)
	case messageTypeCardVote:
		return true, b.voteCard(msg)
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
func (b *Board) boardStateMessage() message {
	// list columns
	columns, err := b.db.ListColumn(b.ID)
	if err != nil {
		log.Printf("broadcastBoardState failed while fetching columns: %s", err)
		columns = []*storage.Column{}
	}

	// list cards
	cards, err := b.db.ListCard(b.ID)
	if err != nil {
		log.Printf("broadcastBoardState failed while fetching cards: %s", err)
		cards = []*storage.Card{}
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

// send sends messages to specific client
func (b *Board) send(msgs []message, c *Client) {
	if c == nil {
		return
	}
	for _, m := range msgs {
		c.message <- m
	}
}

// broadcast sends messages to all clients
func (b *Board) broadcast(msgs []message, exclude *Client) {
	for c := range b.clients {
		if exclude != nil && c.User.ID == exclude.User.ID {
			continue
		}
		for _, m := range msgs {
			c.message <- m
		}
	}
}

// getOrCreateBoard returns a board instance by ID, if not exist, it will create a new board
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
		avatars: make(map[int]bool),
		join:    make(chan *Client),
		leave:   make(chan *Client),
		message: make(chan message),
		stop:    make(chan struct{}),
		timer:   newTimer(),
	}, nil
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
