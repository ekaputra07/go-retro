package model

type MessageType string

const (
	MessageTypeBoardStatus   MessageType = "board.status"
	MessageTypeColumnCreated MessageType = "column.created"
	MessageTypeColumnUpdated MessageType = "column.updated"
	MessageTypeColumnDeleted MessageType = "column.deleted"
	MessageTypeCardNew       MessageType = "card.new"
	MessageTypeCardUpdate    MessageType = "card.update"
	MessageTypeCardDelete    MessageType = "card.delete"
	MessageTypeTimerStarted  MessageType = "timer.started"
	MessageTypeTimerStopped  MessageType = "timer.stopped"
)

type Message struct {
	Type MessageType `json:"type"`
	Data any         `json:"data"`
}
