package board

import (
	"fmt"

	"github.com/google/uuid"
)

// messageType represents the type of message that can be sent to and from the client
type messageType string

const (
	messageTypeBoardUsers        messageType = "board.users"
	messageTypeBoardStatus       messageType = "board.status"
	messageTypeBoardNotification messageType = "board.notification"
	messageTypeColumnNew         messageType = "column.new"
	messageTypeColumnUpdate      messageType = "column.update"
	messageTypeColumnDelete      messageType = "column.delete"
	messageTypeCardNew           messageType = "card.new"
	messageTypeCardUpdate        messageType = "card.update"
	messageTypeCardDelete        messageType = "card.delete"
	messageTypeCardVote          messageType = "card.vote"
	messageTypeTimerCmd          messageType = "timer.cmd"
	messageTypeTimerState        messageType = "timer.state"
)

// message represents a message that can be sent to and from the client
type message struct {
	Type messageType `json:"type"`
	Data any         `json:"data"`

	// from client: assigned when the message is received
	fromClient *Client
}

// dataGet return value(any) from Data by given key
func (m message) dataGet(key string) (any, error) {
	data, ok := m.Data.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("can't convert %+v to map[string]any", m.Data)
	}
	val, ok := data[key]
	if !ok {
		return nil, fmt.Errorf("no key `%s` in %+v", key, data)
	}
	return val, nil
}

// stringVar get string value from data
func (m message) stringVar(to *string, key string) error {
	val, err := m.dataGet(key)
	if err != nil {
		return err
	}
	s, ok := val.(string)
	if !ok {
		return fmt.Errorf("couldn't convert %+v to string", val)
	}
	*to = s
	return nil
}

// intVar get int value from data
func (m message) intVar(to *int, key string) error {
	val, err := m.dataGet(key)
	if err != nil {
		return err
	}
	f, ok := val.(float64)
	if !ok {
		return fmt.Errorf("couldn't convert %+v to float64", val)
	}
	*to = int(f)
	return nil
}

// uuidVar get UUID value from data
func (m message) uuidVar(to *uuid.UUID, key string) error {
	val, err := m.dataGet(key)
	if err != nil {
		return err
	}
	u, err := uuid.Parse(val.(string))
	if err != nil {
		return err
	}
	*to = u
	return nil
}
