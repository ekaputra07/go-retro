package memory

import (
	"fmt"
	"sync"

	"github.com/ekaputra07/go-retro/internal/models"
	"github.com/google/uuid"
)

type columns struct {
	sync.Map

	mu              sync.Mutex
	columnsMaxOrder int
}

func (c *columns) List(boardID uuid.UUID) ([]*models.Column, error) {
	var columns []*models.Column
	c.Range(func(_, v any) bool {
		col := v.(*models.Column)
		if col.BoardID == boardID {
			columns = append(columns, col)
		}
		return true
	})

	return columns, nil
}

func (c *columns) Create(name string, boardID uuid.UUID) (*models.Column, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// auto-increment ordering
	c.columnsMaxOrder++
	nc := models.NewColumn(name, c.columnsMaxOrder, boardID)
	c.Store(nc.ID, nc)
	return nc, nil
}

func (c *columns) Get(id uuid.UUID) (*models.Column, error) {
	if v, ok := c.Load(id); ok {
		return v.(*models.Column), nil
	}
	return nil, fmt.Errorf("column with id=%s doesn't exist", id)
}

func (c *columns) Update(column *models.Column) error {
	if _, loaded := c.LoadOrStore(column.ID, column); !loaded {
		return fmt.Errorf("card with id=%s doesn't exist", column.ID)
	}
	return nil
}

func (c *columns) Delete(id uuid.UUID) error {
	if _, loaded := c.LoadAndDelete(id); !loaded {
		return fmt.Errorf("column with id=%s doesn't exist", id)
	}
	return nil
}
