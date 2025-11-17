package board

import (
	"context"

	"github.com/ekaputra07/go-retro/internal/models"
	"github.com/google/uuid"
)

func (b *Board) createColumn(ctx context.Context, msg message) error {
	var name string
	if err := msg.stringVar(&name, "name"); err != nil {
		return err
	}
	col := models.NewColumn(name, b.store.Columns.NextOrder(), b.ID)
	return b.store.Columns.Create(ctx, col)
}

func (b *Board) deleteColumn(ctx context.Context, msg message) error {
	var id uuid.UUID
	if err := msg.uuidVar(&id, "id"); err != nil {
		return err
	}
	return b.store.Columns.Delete(ctx, id)
}

func (b *Board) updateColumn(ctx context.Context, msg message) error {
	var id uuid.UUID
	if err := msg.uuidVar(&id, "id"); err != nil {
		return err
	}

	col, err := b.store.Columns.Get(ctx, id)
	if err != nil {
		return err
	}

	// if name set and new name is diff, update!
	var name string
	if err := msg.stringVar(&name, "name"); err == nil {
		if name != col.Name {
			col.Name = name
			return b.store.Columns.Update(ctx, *col)
		}
	}
	return nil
}
