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

type cards struct {
	kv jetstream.KeyValue
}

func (u *cards) key(boardID, id uuid.UUID) string {
	return fmt.Sprintf("boards.%s.cards.%s", boardID, id)
}

func (c *cards) List(ctx context.Context, boardID uuid.UUID, limit int) ([]models.Card, error) {
	var cards []models.Card
	lister, err := c.kv.ListKeysFiltered(ctx, fmt.Sprintf("boards.%s.cards.*", boardID))
	if err != nil {
		return cards, err
	}

	counter := 0
	for key := range lister.Keys() {
		val, err := c.kv.Get(ctx, key)
		if err != nil {
			continue // skip
		}
		var c models.Card
		if err = json.Unmarshal(val.Value(), &c); err != nil {
			continue // skip
		}
		cards = append(cards, c)
		counter++
		if counter >= limit {
			lister.Stop()
		}
	}
	return cards, nil
}

func (c *cards) Create(ctx context.Context, card models.Card) error {
	key := c.key(card.BoardID, card.ID)
	_, err := c.kv.Get(ctx, key)
	if err != nil && !errors.Is(err, jetstream.ErrKeyNotFound) {
		return err
	}
	val, err := json.Marshal(card)
	if err != nil {
		return err
	}
	_, err = c.kv.Put(ctx, key, val)
	return err
}

func (c *cards) Get(ctx context.Context, boardID, id uuid.UUID) (*models.Card, error) {
	key := c.key(boardID, id)
	val, err := c.kv.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	var card models.Card
	err = json.Unmarshal(val.Value(), &card)
	return &card, err
}

func (c *cards) Update(ctx context.Context, card models.Card) error {
	b, err := json.Marshal(card)
	if err != nil {
		return err
	}

	_, err = c.kv.Put(ctx, c.key(card.BoardID, card.ID), b)
	return err
}

func (c *cards) Delete(ctx context.Context, boardID, id uuid.UUID) error {
	return c.kv.Delete(ctx, c.key(boardID, id))
}
