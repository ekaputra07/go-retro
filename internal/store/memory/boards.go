package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/ekaputra07/go-retro/internal/models"
	"github.com/google/uuid"
)

type boards struct {
	sync.Map
}

func (b *boards) List(_ context.Context) ([]models.Board, error) {
	var boards []models.Board
	b.Range(func(_, v any) bool {
		board := v.(models.Board)
		boards = append(boards, board)
		return true
	})
	return boards, nil
}

func (b *boards) Create(_ context.Context, board models.Board) error {
	if _, ok := b.Load(board.ID); ok {
		return fmt.Errorf("board with id=%s already exist", board.ID)
	}

	b.Store(board.ID, board)
	return nil
}

func (b *boards) Get(_ context.Context, id uuid.UUID) (*models.Board, error) {
	if v, ok := b.Load(id); ok {
		board := v.(models.Board)
		return &board, nil
	}
	return nil, fmt.Errorf("board with id=%s doesn't exist", id)
}

func (b *boards) Delete(_ context.Context, id uuid.UUID) error {
	if _, loaded := b.LoadAndDelete(id); !loaded {
		return fmt.Errorf("board with id=%s doesn't exist", id)
	}
	return nil
}
