package board

import (
	"errors"

	"github.com/google/uuid"
)

func (b *Board) createColumn(msg *message) error {
	data := msg.Data.(map[string]any)
	name, ok := data["name"]
	if !ok {
		return errors.New("createColumn payload missing `name` field")
	}
	_, err := b.db.CreateColumn(name.(string), b.ID)
	return err
}

func (b *Board) deleteColumn(msg *message) error {
	data := msg.Data.(map[string]any)
	id, ok := data["id"]
	if !ok {
		return errors.New("deleteColumn payload missing `id` field")
	}
	return b.db.DeleteColumn(uuid.MustParse(id.(string)))
}

func (b *Board) updateColumn(msg *message) error {
	data := msg.Data.(map[string]any)

	// get column
	id, ok := data["id"]
	if !ok {
		return errors.New("updateColumn payload missing `id` field")
	}
	col, err := b.db.GetColumn(uuid.MustParse(id.(string)))
	if err != nil {
		return err
	}

	// if name set, update
	name, ok := data["name"]
	if ok {
		col.Name = name.(string)
		return b.db.UpdateColumn(col)
	}
	return nil
}
