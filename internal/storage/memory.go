package storage

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type MemoryStore struct {
	users   sync.Map
	boards  sync.Map
	columns sync.Map
	cards   sync.Map

	mu              sync.Mutex
	columnsMaxOrder int
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		users:           sync.Map{},
		boards:          sync.Map{},
		columns:         sync.Map{},
		cards:           sync.Map{},
		columnsMaxOrder: 0,
	}
}

func (m *MemoryStore) CreateUser() (*User, error) {
	u := NewUser()
	m.users.Store(u.ID, u)
	return u, nil
}

func (m *MemoryStore) GetUser(id uuid.UUID) (*User, error) {
	u, ok := m.users.Load(id)
	if !ok {
		return nil, fmt.Errorf("user id=%s does not exist", id)
	}
	return u.(*User), nil
}

func (m *MemoryStore) UpdateUser(user *User) error {
	_, ok := m.users.Load(user.ID)
	if !ok {
		return fmt.Errorf("user id=%s does not exist", user.ID)
	}
	m.users.Store(user.ID, user)
	return nil
}

func (m *MemoryStore) ListBoard() ([]*Board, error) {
	var boards []*Board
	m.boards.Range(func(_, b any) bool {
		board := b.(*Board)
		boards = append(boards, board)
		return true
	})
	return boards, nil
}

func (m *MemoryStore) CreateBoard(id uuid.UUID) (*Board, error) {
	if _, ok := m.boards.Load(id); ok {
		return nil, fmt.Errorf("board with id=%s already exist", id)
	}

	b := NewBoard(id)
	m.boards.Store(b.ID, b)
	return b, nil
}

func (m *MemoryStore) GetBoard(id uuid.UUID) (*Board, error) {
	if b, ok := m.boards.Load(id); ok {
		return b.(*Board), nil
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

func (m *MemoryStore) ListColumn(boardID uuid.UUID) ([]*Column, error) {
	var columns []*Column
	m.columns.Range(func(_, c any) bool {
		col := c.(*Column)
		if col.BoardID == boardID {
			columns = append(columns, col)
		}
		return true
	})

	return columns, nil
}

func (m *MemoryStore) CreateColumn(name string, boardID uuid.UUID) (*Column, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// auto-increment ordering
	m.columnsMaxOrder++
	c := NewColumn(name, m.columnsMaxOrder, boardID)
	m.columns.Store(c.ID, c)
	return c, nil
}

func (m *MemoryStore) GetColumn(id uuid.UUID) (*Column, error) {
	if c, ok := m.columns.Load(id); ok {
		return c.(*Column), nil
	}
	return nil, fmt.Errorf("column with id=%s doesn't exist", id)
}

func (m *MemoryStore) UpdateColumn(column *Column) error {
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

func (m *MemoryStore) ListCard(boardID uuid.UUID) ([]*Card, error) {
	var cards []*Card
	m.cards.Range(func(_, c any) bool {
		card := c.(*Card)
		if card.BoardID == boardID {
			cards = append(cards, card)
		}
		return true
	})

	return cards, nil
}

func (m *MemoryStore) CreateCard(name string, boardID, columnID uuid.UUID) (*Card, error) {
	c := NewCard(name, boardID, columnID)
	m.cards.Store(c.ID, c)
	return c, nil
}

func (m *MemoryStore) GetCard(id uuid.UUID) (*Card, error) {
	if c, ok := m.cards.Load(id); ok {
		return c.(*Card), nil
	}
	return nil, fmt.Errorf("card with id=%s doesn't exist", id)
}

func (m *MemoryStore) UpdateCard(card *Card) error {
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
