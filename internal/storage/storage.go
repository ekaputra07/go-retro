package storage

import (
	"github.com/google/uuid"
)

type Storage interface {
	// user
	CreateUser() (*User, error)
	GetUser(id uuid.UUID) (*User, error)
	UpdateUser(user *User) error

	// board
	ListBoard() ([]*Board, error)
	CreateBoard(id uuid.UUID) (*Board, error)
	GetBoard(id uuid.UUID) (*Board, error)
	DeleteBoard(id uuid.UUID) error

	// column
	ListColumn(boardID uuid.UUID) ([]*Column, error)
	CreateColumn(name string, boardID uuid.UUID) (*Column, error)
	GetColumn(id uuid.UUID) (*Column, error)
	UpdateColumn(column *Column) error
	DeleteColumn(id uuid.UUID) error

	// card
	ListCard(boardID uuid.UUID) ([]*Card, error)
	CreateCard(name string, boardID, columnID uuid.UUID) (*Card, error)
	GetCard(id uuid.UUID) (*Card, error)
	UpdateCard(card *Card) error
	DeleteCard(id uuid.UUID) error
}
