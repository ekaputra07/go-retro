package board

import (
	"errors"

	"github.com/google/uuid"
)

func (b *Board) createColumn(msg message) error {
	data := msg.Data.(map[string]any)
	name, ok := data["name"]
	if !ok {
		return errors.New("createColumn payload missing `name` field")
	}
	_, err := b.store.Columns.Create(name.(string), b.ID)
	return err
}

func (b *Board) deleteColumn(msg message) error {
	data := msg.Data.(map[string]any)
	id, ok := data["id"]
	if !ok {
		return errors.New("deleteColumn payload missing `id` field")
	}
	return b.store.Columns.Delete(uuid.MustParse(id.(string)))
}

func (b *Board) updateColumn(msg message) error {
	data := msg.Data.(map[string]any)

	// get column
	id, ok := data["id"]
	if !ok {
		return errors.New("updateColumn payload missing `id` field")
	}
	col, err := b.store.Columns.Get(uuid.MustParse(id.(string)))
	if err != nil {
		return err
	}

	// if name set, update
	name, ok := data["name"]
	if ok {
		col.Name = name.(string)
		return b.store.Columns.Update(col)
	}
	return nil
}
