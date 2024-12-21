package server

import (
	"time"

	"github.com/ekaputra07/go-retro/internal/model"
	"github.com/google/uuid"
)

type card struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Column    uuid.UUID `json:"column"`
	CreatedAt int64     `json:"created_at"`
}

func (b *board) createCard(msg *model.Message) {
	// TODO: need a cleaner way to handle this
	data := msg.Data.(map[string]any)
	name, ok := data["name"]
	if !ok {
		return
	}
	column, ok := data["column"]
	if !ok {
		return
	}
	card := &card{
		ID:        uuid.New(),
		Name:      name.(string),
		Column:    uuid.MustParse(column.(string)),
		CreatedAt: time.Now().Unix(),
	}
	b.Cards = append(b.Cards, card)
	b.broadcastStatus()
}

func (b *board) deleteCard(msg *model.Message) {
	data := msg.Data.(map[string]any)
	id, ok := data["id"]
	if !ok {
		return
	}
	cards := []*card{}
	for _, c := range b.Cards {
		if c.ID.String() != id {
			cards = append(cards, c)
		}
	}
	b.Cards = cards
	b.broadcastStatus()
}

func (b *board) updateCard(msg *model.Message) {
	data := msg.Data.(map[string]any)
	id, ok := data["id"]
	if !ok {
		return
	}
	name, ok := data["name"]
	if !ok {
		return
	}
	column, ok := data["column"]
	if !ok {
		return
	}
	card := b.cardById(uuid.MustParse(id.(string)))
	if card == nil {
		return
	}
	card.Name = name.(string)
	card.Column = uuid.MustParse(column.(string))
	b.broadcastStatus()
}

func (b *board) cardById(id uuid.UUID) *card {
	for _, c := range b.Cards {
		if c.ID == id {
			return c
		}
	}
	return nil
}
