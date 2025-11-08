package memory

import (
	"fmt"
	"sync"

	"github.com/ekaputra07/go-retro/internal/models"
	"github.com/google/uuid"
)

type boards struct {
	sync.Map
}

func (b *boards) List() ([]*models.Board, error) {
	var boards []*models.Board
	b.Range(func(_, v any) bool {
		board := v.(*models.Board)
		boards = append(boards, board)
		return true
	})
	return boards, nil
}

func (b *boards) Create(id uuid.UUID) (*models.Board, error) {
	if _, ok := b.Load(id); ok {
		return nil, fmt.Errorf("board with id=%s already exist", id)
	}

	nb := models.NewBoard(id)
	b.Store(nb.ID, nb)
	return nb, nil
}

func (b *boards) Get(id uuid.UUID) (*models.Board, error) {
	if v, ok := b.Load(id); ok {
		return v.(*models.Board), nil
	}
	return nil, fmt.Errorf("board with id=%s doesn't exist", id)
}

func (b *boards) Delete(id uuid.UUID) error {
	if _, loaded := b.LoadAndDelete(id); !loaded {
		return fmt.Errorf("board with id=%s doesn't exist", id)
	}
	return nil
}
