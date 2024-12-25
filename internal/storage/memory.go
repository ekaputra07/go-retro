package storage

import (
	"fmt"
	"sync"

	"github.com/ekaputra07/go-retro/internal/model"
	"github.com/google/uuid"
)

type MemoryStore struct {
	boards  sync.Map
	columns sync.Map
	cards   sync.Map

	mu              sync.Mutex
	columnsMaxOrder int
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		boards:          sync.Map{},
		columns:         sync.Map{},
		cards:           sync.Map{},
		columnsMaxOrder: 0,
	}
}

func (m *MemoryStore) ListBoard() ([]*model.Board, error) {
	var boards []*model.Board
	m.boards.Range(func(_, b any) bool {
		board := b.(*model.Board)
		boards = append(boards, board)
		return true
	})
	return boards, nil
}

func (m *MemoryStore) CreateBoard(id uuid.UUID) (*model.Board, error) {
	if _, ok := m.boards.Load(id); ok {
		return nil, fmt.Errorf("board with id=%s already exist", id)
	}

	b := model.NewBoard(id)
	m.boards.Store(b.ID, b)
	return b, nil
}

func (m *MemoryStore) GetBoard(id uuid.UUID) (*model.Board, error) {
	if b, ok := m.boards.Load(id); ok {
		return b.(*model.Board), nil
	}
	return nil, fmt.Errorf("board with id=%s doesn't exist", id)
}

func (m *MemoryStore) DeleteBoard(id uuid.UUID) error {
	if _, ok := m.boards.Load(id); ok {
		m.boards.Delete(id)
		return nil
	}
	return fmt.Errorf("board with id=%s doesn't exist", id)
}

func (m *MemoryStore) ListColumn(boardID uuid.UUID) ([]*model.Column, error) {
	var columns []*model.Column
	m.columns.Range(func(_, c any) bool {
		col := c.(*model.Column)
		if col.BoardID == boardID {
			columns = append(columns, col)
		}
		return true
	})

	return columns, nil
}

func (m *MemoryStore) CreateColumn(name string, boardID uuid.UUID) (*model.Column, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// auto-increment ordering
	m.columnsMaxOrder++
	c := model.NewColumn(name, m.columnsMaxOrder, boardID)
	m.columns.Store(c.ID, c)
	return c, nil
}

func (m *MemoryStore) GetColumn(id uuid.UUID) (*model.Column, error) {
	if c, ok := m.columns.Load(id); ok {
		return c.(*model.Column), nil
	}
	return nil, fmt.Errorf("column with id=%s doesn't exist", id)
}

func (m *MemoryStore) UpdateColumn(column *model.Column) error {
	m.columns.Store(column.ID, column)
	return nil
}

func (m *MemoryStore) DeleteColumn(id uuid.UUID) error {
	if _, ok := m.columns.Load(id); ok {
		m.columns.Delete(id)
		return nil
	}
	return fmt.Errorf("column with id=%s doesn't exist", id)
}

func (m *MemoryStore) ListCard(boardID uuid.UUID) ([]*model.Card, error) {
	var cards []*model.Card
	m.cards.Range(func(_, c any) bool {
		card := c.(*model.Card)
		if card.BoardID == boardID {
			cards = append(cards, card)
		}
		return true
	})

	return cards, nil
}

func (m *MemoryStore) CreateCard(name string, boardID, columnID uuid.UUID) (*model.Card, error) {
	c := model.NewCard(name, boardID, columnID)
	m.cards.Store(c.ID, c)
	return c, nil
}

func (m *MemoryStore) GetCard(id uuid.UUID) (*model.Card, error) {
	if c, ok := m.cards.Load(id); ok {
		return c.(*model.Card), nil
	}
	return nil, fmt.Errorf("card with id=%s doesn't exist", id)
}

func (m *MemoryStore) UpdateCard(card *model.Card) error {
	m.cards.Store(card.ID, card)
	return nil
}

func (m *MemoryStore) DeleteCard(id uuid.UUID) error {
	if _, ok := m.cards.Load(id); ok {
		m.cards.Delete(id)
		return nil
	}
	return fmt.Errorf("card with id=%s doesn't exist", id)
}
