package board

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ekaputra07/go-retro/internal/models"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go/jetstream"
)

// stream represent a single item from a stream of changes
type stream struct {
	Type   string `json:"type"`
	ID     string `json:"id"`
	Op     string `json:"op"`
	Object any    `json:"obj"`
}

func newStreamFromKVE(kve jetstream.KeyValueEntry) (*stream, error) {
	// key format: boards.<id>.<type>.<id>
	tokens := strings.Split(kve.Key(), ".")
	s := stream{Type: tokens[2], ID: tokens[3]}

	switch kve.Operation() {
	case jetstream.KeyValuePut:
		s.Op = "put"
	case jetstream.KeyValueDelete:
		s.Op = "del"
		return &s, nil // return without Object
	case jetstream.KeyValuePurge:
		s.Op = "del" // return without Object
		return &s, nil
	}

	switch s.Type {
	case "columns":
		var c models.Column
		if err := json.Unmarshal(kve.Value(), &c); err != nil {
			return nil, err
		}
		s.Object = c
	case "cards":
		var c models.Card
		if err := json.Unmarshal(kve.Value(), &c); err != nil {
			return nil, err
		}
		s.Object = c
	default:
		return nil, fmt.Errorf("stream type %s not supported", s.Type)
	}
	return &s, nil
}

// messageType represents the type of message that can be sent to and from the client
type messageType string

const (
	messageTypeMe                messageType = "me"
	messageTypeBoardUsers        messageType = "board.users"
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
	User models.User `json:"user"`
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
