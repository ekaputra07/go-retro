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
	c := &card{
		ID:        uuid.New(),
		Name:      name.(string),
		Column:    uuid.MustParse(column.(string)),
		CreatedAt: time.Now().Unix(),
	}
	b.Cards = append(b.Cards, c)
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
	c := b.cardById(uuid.MustParse(id.(string)))
	if c == nil {
		return
	}
	c.Name = name.(string)
	c.Column = uuid.MustParse(column.(string))
}

func (b *board) cardById(id uuid.UUID) *card {
	for _, c := range b.Cards {
		if c.ID == id {
			return c
		}
	}
	return nil
}
