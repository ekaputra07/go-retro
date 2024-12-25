package model

import (
	"time"

	"github.com/google/uuid"
)

type MessageType string

const (
	MessageTypeBoardStatus  MessageType = "board.status"
	MessageTypeColumnNew    MessageType = "column.new"
	MessageTypeColumnUpdate MessageType = "column.update"
	MessageTypeColumnDelete MessageType = "column.delete"
	MessageTypeCardNew      MessageType = "card.new"
	MessageTypeCardUpdate   MessageType = "card.update"
	MessageTypeCardDelete   MessageType = "card.delete"
	MessageTypeTimerStarted MessageType = "timer.started"
	MessageTypeTimerStopped MessageType = "timer.stopped"
)

type Message struct {
	Type MessageType `json:"type"`
	Data any         `json:"data"`
}

type Board struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt int64     `json:"created_at"`
}

func NewBoard(id uuid.UUID) *Board {
	return &Board{
		ID:        id,
		CreatedAt: time.Now().Unix(),
	}
}

type Column struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Order     int       `json:"order"`
	BoardID   uuid.UUID `json:"board_id"`
	CreatedAt int64     `json:"created_at"`
}

func NewColumn(name string, order int, boardID uuid.UUID) *Column {
	return &Column{
		ID:        uuid.New(),
		Name:      name,
		BoardID:   boardID,
		Order:     order,
		CreatedAt: time.Now().Unix(),
	}
}

type Card struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	BoardID   uuid.UUID `json:"board_id"`
	ColumnID  uuid.UUID `json:"column_id"`
	CreatedAt int64     `json:"created_at"`
}

func NewCard(name string, boardID, columnID uuid.UUID) *Card {
	return &Card{
		ID:        uuid.New(),
		Name:      name,
		BoardID:   boardID,
		ColumnID:  columnID,
		CreatedAt: time.Now().Unix(),
	}
}
