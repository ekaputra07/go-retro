package natstore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ekaputra07/go-retro/internal/models"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go/jetstream"
)

type columns struct {
	kv jetstream.KeyValue
}

func (u *columns) key(boardID, id uuid.UUID) string {
	return fmt.Sprintf("boards.%s.columns.%s", boardID, id)
}

func (c *columns) List(ctx context.Context, boardID uuid.UUID) ([]models.Column, error) {
	var columns []models.Column
	lister, err := c.kv.ListKeysFiltered(ctx, fmt.Sprintf("boards.%s.columns.*", boardID))
	if err != nil {
		return columns, err
	}
	for key := range lister.Keys() {
		val, err := c.kv.Get(ctx, key)
		if err != nil {
			continue // skip
		}
		var c models.Column
		if err = json.Unmarshal(val.Value(), &c); err != nil {
			continue // skip
		}
		columns = append(columns, c)
	}
	return columns, nil
}

func (c *columns) Create(ctx context.Context, column models.Column) error {
	key := c.key(column.BoardID, column.ID)
	_, err := c.kv.Get(ctx, key)
	if err != nil && !errors.Is(err, jetstream.ErrKeyNotFound) {
		return err
	}
	val, err := json.Marshal(column)
	if err != nil {
		return err
	}
	_, err = c.kv.Put(ctx, key, val)
	return err
}

func (c *columns) Get(ctx context.Context, boardID, id uuid.UUID) (*models.Column, error) {
	key := c.key(boardID, id)
	val, err := c.kv.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	var column models.Column
	err = json.Unmarshal(val.Value(), &column)
	return &column, err
}

func (c *columns) Update(ctx context.Context, column models.Column) error {
	b, err := json.Marshal(column)
	if err != nil {
		return err
	}

	_, err = c.kv.Put(ctx, c.key(column.BoardID, column.ID), b)
	return err
}

func (c *columns) Delete(ctx context.Context, boardID, id uuid.UUID) error {
	return c.kv.Delete(ctx, c.key(boardID, id))
}
