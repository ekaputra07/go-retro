package memstore

import (
	"context"
	"testing"

	"github.com/ekaputra07/go-retro/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	c := &cards{}
	cards, _ := c.List(context.Background())
	assert.Len(t, cards, 0)
}

func TestCreateCard(t *testing.T) {
	ctx := context.Background()

	c := &cards{}
	columnId := uuid.New()
	c.Create(ctx, models.NewCard("card1", boardID, columnId))
	c.Create(ctx, models.NewCard("card2", boardID, columnId))

	cards, _ := c.List(ctx)
	assert.Len(t, cards, 2)
}

func TestGetCard(t *testing.T) {
	ctx := context.Background()

	c := &cards{}
	card := models.NewCard("card1", boardID, uuid.New())
	c.Create(ctx, card)

	_, err := c.Get(ctx, uuid.New())
	assert.Error(t, err)

	_, err = c.Get(ctx, card.ID)
	assert.NoError(t, err)
}

func TestUpdateCard(t *testing.T) {
	ctx := context.Background()

	c := &cards{}
	card := models.NewCard("card1", boardID, uuid.New())
	c.Create(ctx, card)

	card.Name = "card-1-updated"
	err := c.Update(ctx, card)
	assert.NoError(t, err)

	ca, _ := c.Get(ctx, card.ID)
	assert.Equal(t, "card-1-updated", ca.Name)
}

func TestDeleteCard(t *testing.T) {
	ctx := context.Background()

	c := &cards{}
	columnId := uuid.New()
	card := models.NewCard("card1", boardID, columnId)
	c.Create(ctx, card)
	c.Create(ctx, models.NewCard("card2", boardID, columnId))

	err := c.Delete(ctx, card.ID)
	assert.NoError(t, err)

	cards, _ := c.List(ctx)
	assert.Len(t, cards, 1)
}
