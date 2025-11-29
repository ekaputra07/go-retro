package store

import (
	"context"

	"github.com/ekaputra07/go-retro/internal/models"
	"github.com/google/uuid"
)

type UserRepo interface {
	Create(ctx context.Context, user models.User) error
	Get(ctx context.Context, id uuid.UUID) (*models.User, error)
	Update(ctx context.Context, user models.User) error
}

type BoardRepo interface {
	List(ctx context.Context, limit int) ([]models.Board, error)
	Create(ctx context.Context, board models.Board) error
	Get(ctx context.Context, id uuid.UUID) (*models.Board, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type ColumnRepo interface {
	ListKeys(ctx context.Context, boardID uuid.UUID, limit int) ([]string, error)
	List(ctx context.Context, boardID uuid.UUID, limit int) ([]models.Column, error)
	Create(ctx context.Context, column models.Column) error
	Get(ctx context.Context, boardID uuid.UUID, id uuid.UUID) (*models.Column, error)
	Update(ctx context.Context, column models.Column) error
	Delete(ctx context.Context, boardID uuid.UUID, id uuid.UUID) error
}

type CardRepo interface {
	List(ctx context.Context, boardID uuid.UUID, limit int) ([]models.Card, error)
	Create(ctx context.Context, card models.Card) error
	Get(ctx context.Context, boardID uuid.UUID, id uuid.UUID) (*models.Card, error)
	Update(ctx context.Context, card models.Card) error
	Delete(ctx context.Context, boardID uuid.UUID, id uuid.UUID) error
}

// Store stores globally available records e.g Users and Boards
type Store struct {
	Users   UserRepo
	Boards  BoardRepo
	Columns ColumnRepo
	Cards   CardRepo
}
