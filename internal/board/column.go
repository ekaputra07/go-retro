package board

import (
	"github.com/google/uuid"
)

func (b *Board) createColumn(msg message) error {
	var name string
	if err := msg.stringVar(&name, "name"); err != nil {
		return err
	}
	_, err := b.store.Columns.Create(name, b.ID)
	return err
}

func (b *Board) deleteColumn(msg message) error {
	var id uuid.UUID
	if err := msg.uuidVar(&id, "id"); err != nil {
		return err
	}
	return b.store.Columns.Delete(id)
}

func (b *Board) updateColumn(msg message) error {
	var id uuid.UUID
	if err := msg.uuidVar(&id, "id"); err != nil {
		return err
	}

	col, err := b.store.Columns.Get(id)
	if err != nil {
		return err
	}

	// if name set and new name is diff, update!
	var name string
	if err := msg.stringVar(&name, "name"); err == nil {
		if name != col.Name {
			col.Name = name
			return b.store.Columns.Update(col)
		}
	}
	return nil
}
