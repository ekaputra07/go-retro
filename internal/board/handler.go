package board

import (
	"context"
	"fmt"

	"github.com/ekaputra07/go-retro/internal/models"
	"github.com/ekaputra07/go-retro/internal/store"
	"github.com/google/uuid"
)

// messageHandler handles incoming message and operates on the store.
type messageHandler struct {
	store *store.Store
}

func newMessageHandler(store *store.Store) *messageHandler {
	return &messageHandler{store}
}

func (h *messageHandler) handle(ctx context.Context, msg message) error {
	switch msg.Type {
	case messageTypeColumnNew:
		return h.createColumn(ctx, msg)
	case messageTypeColumnDelete:
		return h.deleteColumn(ctx, msg)
	case messageTypeColumnUpdate:
		return h.updateColumn(ctx, msg)
	case messageTypeCardNew:
		return h.createCard(ctx, msg)
	case messageTypeCardDelete:
		return h.deleteCard(ctx, msg)
	case messageTypeCardUpdate:
		return h.updateCard(ctx, msg)
	case messageTypeCardVote:
		return h.voteCard(ctx, msg)
	}
	return fmt.Errorf("message type=%s not supported by messageHandler", msg.Type)
}

func (h *messageHandler) createColumn(ctx context.Context, msg message) error {
	var name string
	if err := msg.stringVar(&name, "name"); err != nil {
		return err
	}
	col := models.NewColumn(name, msg.BoardID)
	return h.store.Columns.Create(ctx, col)
}

func (h *messageHandler) deleteColumn(ctx context.Context, msg message) error {
	var id uuid.UUID
	if err := msg.uuidVar(&id, "id"); err != nil {
		return err
	}
	return h.store.Columns.Delete(ctx, msg.BoardID, id)
}

func (h *messageHandler) updateColumn(ctx context.Context, msg message) error {
	var id uuid.UUID
	if err := msg.uuidVar(&id, "id"); err != nil {
		return err
	}

	col, err := h.store.Columns.Get(ctx, msg.BoardID, id)
	if err != nil {
		return err
	}

	// if name set and new name is diff, update!
	var name string
	if err := msg.stringVar(&name, "name"); err == nil {
		if name != col.Name {
			col.Name = name
			return h.store.Columns.Update(ctx, *col)
		}
	}
	return nil
}

func (h *messageHandler) createCard(ctx context.Context, msg message) error {
	var name string
	var columnID uuid.UUID

	if err := msg.stringVar(&name, "name"); err != nil {
		return err
	}
	if err := msg.uuidVar(&columnID, "column_id"); err != nil {
		return err
	}
	col, err := h.store.Columns.Get(ctx, msg.BoardID, columnID)
	if err != nil {
		return err
	}
	card := models.NewCard(name, msg.BoardID, col.ID)
	return h.store.Cards.Create(ctx, card)
}

func (h *messageHandler) deleteCard(ctx context.Context, msg message) error {
	var id uuid.UUID
	if err := msg.uuidVar(&id, "id"); err != nil {
		return err
	}
	return h.store.Cards.Delete(ctx, msg.BoardID, id)
}

func (h *messageHandler) updateCard(ctx context.Context, msg message) error {
	var id uuid.UUID
	var name string
	var columnID uuid.UUID

	if err := msg.uuidVar(&id, "id"); err != nil {
		return err
	}

	card, err := h.store.Cards.Get(ctx, msg.BoardID, id)
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
	return h.store.Cards.Update(ctx, *card)
}

func (h *messageHandler) voteCard(ctx context.Context, msg message) error {
	var id uuid.UUID
	var vote int

	if err := msg.uuidVar(&id, "id"); err != nil {
		return err
	}

	card, err := h.store.Cards.Get(ctx, msg.BoardID, id)
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
	return h.store.Cards.Update(ctx, *card)
}
