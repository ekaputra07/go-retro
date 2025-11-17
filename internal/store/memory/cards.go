package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/ekaputra07/go-retro/internal/models"
	"github.com/google/uuid"
)

type cards struct {
	sync.Map
}

func (c *cards) List(_ context.Context) ([]models.Card, error) {
	var cards []models.Card
	c.Range(func(_, v any) bool {
		card := v.(models.Card)
		cards = append(cards, card)
		return true
	})

	return cards, nil
}

func (c *cards) Create(_ context.Context, card models.Card) error {
	c.Store(card.ID, card)
	return nil
}

func (c *cards) Get(_ context.Context, id uuid.UUID) (*models.Card, error) {
	if v, ok := c.Load(id); ok {
		card := v.(models.Card)
		return &card, nil
	}
	return nil, fmt.Errorf("card with id=%s doesn't exist", id)
}

func (c *cards) Update(_ context.Context, card models.Card) error {
	if _, ok := c.Load(card.ID); !ok {
		return fmt.Errorf("card with id=%s doesn't exist", card.ID)
	}
	c.Store(card.ID, card)
	return nil
}

func (c *cards) Delete(_ context.Context, id uuid.UUID) error {
	if _, loaded := c.LoadAndDelete(id); !loaded {
		return fmt.Errorf("card with id=%s doesn't exist", id)
	}
	return nil
}
