package server

import (
	"errors"
	"fmt"

	"github.com/ekaputra07/go-retro/internal/model"
	"github.com/google/uuid"
)

func (b *board) createCard(msg *model.Message) error {
	data := msg.Data.(map[string]any)
	name, ok := data["name"]
	if !ok {
		return errors.New("createCard payload missing `name` field")
	}
	colID, ok := data["column_id"]
	if !ok {
		return errors.New("createCard payload missing `column_id` field")
	}
	col, err := b.db.GetColumn(uuid.MustParse(colID.(string)))
	if err != nil {
		return err
	}
	_, err = b.db.CreateCard(name.(string), b.ID, col.ID)
	return err
}

func (b *board) deleteCard(msg *model.Message) error {
	data := msg.Data.(map[string]any)
	id, ok := data["id"]
	if !ok {
		return errors.New("deleteCard payload missing `id` field")
	}
	return b.db.DeleteCard(uuid.MustParse(id.(string)))
}

func (b *board) updateCard(msg *model.Message) error {
	data := msg.Data.(map[string]any)

	// get card
	id, ok := data["id"]
	if !ok {
		return errors.New("updateCard payload missing `id` field")
	}
	card, err := b.db.GetCard(uuid.MustParse(id.(string)))
	if err != nil {
		return err
	}

	// if name set, update
	name, ok := data["name"]
	if ok {
		card.Name = name.(string)
	}

	// if colum set, update
	colID, ok := data["column_id"]
	if ok {
		col, err := b.db.GetColumn(uuid.MustParse(colID.(string)))
		if err != nil {
			return err
		}
		card.ColumnID = col.ID
	}
	return b.db.UpdateCard(card)
}

func (b *board) voteCard(msg *model.Message) error {
	data := msg.Data.(map[string]any)

	// get card
	id, ok := data["id"]
	if !ok {
		return errors.New("voteCard payload missing `id` field")
	}
	card, err := b.db.GetCard(uuid.MustParse(id.(string)))
	if err != nil {
		return err
	}

	// if vote set, update
	vote, ok := data["vote"]
	if ok {
		v := vote.(int)
		if v != 1 && v != -1 {
			return fmt.Errorf("vote value of %v is invalid", v)
		}
		card.Votes += v
		return b.db.UpdateCard(card)
	}
	return nil
}
