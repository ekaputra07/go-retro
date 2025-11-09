package board

import (
	"fmt"

	"github.com/google/uuid"
)

func (b *Board) createCard(msg message) error {
	var name string
	var columnID uuid.UUID

	if err := msg.stringVar(&name, "name"); err != nil {
		return err
	}
	if err := msg.uuidVar(&columnID, "column_id"); err != nil {
		return err
	}
	col, err := b.store.Columns.Get(columnID)
	if err != nil {
		return err
	}
	_, err = b.store.Cards.Create(name, b.ID, col.ID)
	return err
}

func (b *Board) deleteCard(msg message) error {
	var id uuid.UUID
	if err := msg.uuidVar(&id, "id"); err != nil {
		return err
	}
	return b.store.Cards.Delete(id)
}

func (b *Board) updateCard(msg message) error {
	var id uuid.UUID
	var name string
	var columnID uuid.UUID

	if err := msg.uuidVar(&id, "id"); err != nil {
		return err
	}

	card, err := b.store.Cards.Get(id)
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
	return b.store.Cards.Update(card)
}

func (b *Board) voteCard(msg message) error {
	var id uuid.UUID
	var vote int

	if err := msg.uuidVar(&id, "id"); err != nil {
		return err
	}

	card, err := b.store.Cards.Get(id)
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
	return b.store.Cards.Update(card)
}
