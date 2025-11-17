package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/ekaputra07/go-retro/internal/models"
	"github.com/google/uuid"
)

type columns struct {
	sync.Map

	mu        sync.Mutex
	nextOrder int
}

func (c *columns) NextOrder() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.nextOrder++
	return c.nextOrder
}

func (c *columns) List(_ context.Context) ([]models.Column, error) {
	var columns []models.Column
	c.Range(func(_, v any) bool {
		col := v.(models.Column)
		columns = append(columns, col)
		return true
	})

	return columns, nil
}

func (c *columns) Create(_ context.Context, column models.Column) error {
	c.Store(column.ID, column)
	return nil
}

func (c *columns) Get(_ context.Context, id uuid.UUID) (*models.Column, error) {
	if v, ok := c.Load(id); ok {
		col := v.(models.Column)
		return &col, nil
	}
	return nil, fmt.Errorf("column with id=%s doesn't exist", id)
}

func (c *columns) Update(_ context.Context, column models.Column) error {
	if _, ok := c.Load(column.ID); !ok {
		return fmt.Errorf("card with id=%s doesn't exist", column.ID)
	}
	c.Store(column.ID, column)
	return nil
}

func (c *columns) Delete(_ context.Context, id uuid.UUID) error {
	if _, loaded := c.LoadAndDelete(id); !loaded {
		return fmt.Errorf("column with id=%s doesn't exist", id)
	}
	return nil
}
