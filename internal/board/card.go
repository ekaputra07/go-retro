package board

import (
	"context"
	"fmt"

	"github.com/ekaputra07/go-retro/internal/models"
	"github.com/google/uuid"
)

func (b *Board) createCard(ctx context.Context, msg message) error {
	var name string
	var columnID uuid.UUID

	if err := msg.stringVar(&name, "name"); err != nil {
		return err
	}
	if err := msg.uuidVar(&columnID, "column_id"); err != nil {
		return err
	}
	col, err := b.store.Columns.Get(ctx, columnID)
	if err != nil {
		return err
	}
	card := models.NewCard(name, b.ID, col.ID)
	return b.store.Cards.Create(ctx, card)
}

func (b *Board) deleteCard(ctx context.Context, msg message) error {
	var id uuid.UUID
	if err := msg.uuidVar(&id, "id"); err != nil {
		return err
	}
	return b.store.Cards.Delete(ctx, id)
}

func (b *Board) updateCard(ctx context.Context, msg message) error {
	var id uuid.UUID
	var name string
	var columnID uuid.UUID

	if err := msg.uuidVar(&id, "id"); err != nil {
		return err
	}

	card, err := b.store.Cards.Get(ctx, id)
	if err != nil {
		return err
	}

	// update card name if new name given
	if err := msg.stringVar(&name, "name"); err == nil {
		if name != card.Name {
			card.Name = name
		}
	}
	// move to different column if new column_id given
	if err := msg.uuidVar(&columnID, "column_id"); err == nil {
		if columnID != card.ColumnID {
			card.ColumnID = columnID
		}
	}
	return b.store.Cards.Update(ctx, *card)
}

func (b *Board) voteCard(ctx context.Context, msg message) error {
	var id uuid.UUID
	var vote int

	if err := msg.uuidVar(&id, "id"); err != nil {
		return err
	}

	card, err := b.store.Cards.Get(ctx, id)
	if err != nil {
		return err
	}

	if err := msg.intVar(&vote, "vote"); err != nil {
		return err
	}
	if vote != 1 && vote != -1 {
		return fmt.Errorf("vote value of %v is invalid", vote)
	}

	card.Votes += vote
	return b.store.Cards.Update(ctx, *card)
}
