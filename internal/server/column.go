package server

import (
	"time"

	"github.com/ekaputra07/go-retro/internal/model"
	"github.com/google/uuid"
)

var defaultColumns = []string{"Good", "Bad", "Questions"}

type column struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt int64     `json:"created_at"`
}

func (b *board) columnByID(id uuid.UUID) *column {
	for _, c := range b.Columns {
		if c.ID == id {
			return c
		}
	}
	return nil
}

func (b *board) createColumn(msg *model.Message) {
	data := msg.Data.(map[string]any)
	name, ok := data["name"]
	if !ok {
		return
	}
	c := &column{
		ID:        uuid.New(),
		Name:      name.(string),
		CreatedAt: time.Now().Unix(),
	}
	b.Columns = append(b.Columns, c)
}

func (b *board) deleteColumn(msg *model.Message) {
	data := msg.Data.(map[string]any)
	id, ok := data["id"]
	if !ok {
		return
	}
	columns := []*column{}
	for _, c := range b.Columns {
		if c.ID.String() != id {
			columns = append(columns, c)
		}
	}
	b.Columns = columns
}

func (b *board) updateColumn(msg *model.Message) {
	data := msg.Data.(map[string]any)
	id, ok := data["id"]
	if !ok {
		return
	}
	name, ok := data["name"]
	if !ok {
		return
	}
	c := b.columnByID(uuid.MustParse(id.(string)))
	if c == nil {
		return
	}
	c.Name = name.(string)
}
