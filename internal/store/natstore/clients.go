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

type clients struct {
	kv jetstream.KeyValue
}

func (c *clients) key(boardID, id uuid.UUID) string {
	return fmt.Sprintf("boards.%s.clients.%s", boardID, id)
}

func (c *clients) Create(ctx context.Context, client models.Client) error {
	key := c.key(client.BoardID, client.ID)
	_, err := c.kv.Get(ctx, key)
	if err != nil && !errors.Is(err, jetstream.ErrKeyNotFound) {
		return err
	}
	val, err := json.Marshal(client)
	if err != nil {
		return err
	}
	_, err = c.kv.Put(ctx, key, val)
	return err
}

func (c *clients) Delete(ctx context.Context, boardID, id uuid.UUID) error {
	return c.kv.Delete(ctx, c.key(boardID, id))
}
