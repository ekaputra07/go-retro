package store

import (
	"github.com/ekaputra07/go-retro/internal/models"
	"github.com/google/uuid"
)

type UserStore interface {
	Create(avatarID int) (*models.User, error)
	Get(id uuid.UUID) (*models.User, error)
	Update(user *models.User) error
}

type BoardStore interface {
	List() ([]*models.Board, error)
	Create(id uuid.UUID) (*models.Board, error)
	Get(id uuid.UUID) (*models.Board, error)
	Delete(id uuid.UUID) error
}

type ColumnStore interface {
	List(boardID uuid.UUID) ([]*models.Column, error)
	Create(name string, boardID uuid.UUID) (*models.Column, error)
	Get(id uuid.UUID) (*models.Column, error)
	Update(column *models.Column) error
	Delete(id uuid.UUID) error
}

type CardStore interface {
	List(boardID uuid.UUID) ([]*models.Card, error)
	Create(name string, boardID, columnID uuid.UUID) (*models.Card, error)
	Get(id uuid.UUID) (*models.Card, error)
	Update(card *models.Card) error
	Delete(id uuid.UUID) error
}

type Store struct {
	Users   UserStore
	Boards  BoardStore
	Columns ColumnStore
	Cards   CardStore
}
