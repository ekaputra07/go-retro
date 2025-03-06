package storage

import (
	"github.com/ekaputra07/go-retro/internal/model"
	"github.com/google/uuid"
)

type Storage interface {
	// user
	CreateUser() (*model.User, error)
	GetUser(id uuid.UUID) (*model.User, error)
	UpdateUser(user *model.User) error

	// board
	ListBoard() ([]*model.Board, error)
	CreateBoard(id uuid.UUID) (*model.Board, error)
	GetBoard(id uuid.UUID) (*model.Board, error)
	DeleteBoard(id uuid.UUID) error

	// column
	ListColumn(boardID uuid.UUID) ([]*model.Column, error)
	CreateColumn(name string, boardID uuid.UUID) (*model.Column, error)
	GetColumn(id uuid.UUID) (*model.Column, error)
	UpdateColumn(column *model.Column) error
	DeleteColumn(id uuid.UUID) error

	// card
	ListCard(boardID uuid.UUID) ([]*model.Card, error)
	CreateCard(name string, boardID, columnID uuid.UUID) (*model.Card, error)
	GetCard(id uuid.UUID) (*model.Card, error)
	UpdateCard(card *model.Card) error
	DeleteCard(id uuid.UUID) error
}
