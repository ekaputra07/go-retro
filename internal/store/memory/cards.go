package memory

import (
	"fmt"
	"sync"

	"github.com/ekaputra07/go-retro/internal/models"
	"github.com/google/uuid"
)

type cards struct {
	sync.Map
}

func (c *cards) List(boardID uuid.UUID) ([]*models.Card, error) {
	var cards []*models.Card
	c.Range(func(_, v any) bool {
		card := v.(*models.Card)
		if card.BoardID == boardID {
			cards = append(cards, card)
		}
		return true
	})

	return cards, nil
}

func (c *cards) Create(name string, boardID, columnID uuid.UUID) (*models.Card, error) {
	nc := models.NewCard(name, boardID, columnID)
	c.Store(nc.ID, nc)
	return nc, nil
}

func (c *cards) Get(id uuid.UUID) (*models.Card, error) {
	if v, ok := c.Load(id); ok {
		return v.(*models.Card), nil
	}
	return nil, fmt.Errorf("card with id=%s doesn't exist", id)
}

func (c *cards) Update(card *models.Card) error {
	if _, loaded := c.LoadOrStore(card.ID, card); !loaded {
		return fmt.Errorf("card with id=%s doesn't exist", card.ID)
	}
	return nil
}

func (c *cards) Delete(id uuid.UUID) error {
	if _, loaded := c.LoadAndDelete(id); !loaded {
		return fmt.Errorf("card with id=%s doesn't exist", id)
	}
	return nil
}
