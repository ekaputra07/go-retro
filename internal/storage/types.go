package storage

import (
	"time"

	"github.com/google/uuid"
)

// User holds information the person that join the board.
// Users are not bounded to specific board but are global entities.
// Allowed to join multiple boards OR join a single board through multiple connection (client)
type User struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func NewUser() *User {
	return &User{
		ID: uuid.New(),
	}
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
	Votes     int       `json:"votes"`
	CreatedAt int64     `json:"created_at"`
}

func NewCard(name string, boardID, columnID uuid.UUID) *Card {
	return &Card{
		ID:        uuid.New(),
		Name:      name,
		BoardID:   boardID,
		ColumnID:  columnID,
		Votes:     0,
		CreatedAt: time.Now().Unix(),
	}
}
