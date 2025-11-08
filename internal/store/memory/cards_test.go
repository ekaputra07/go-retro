package memory

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	c := &cards{}
	cards, _ := c.List(boardID)
	assert.Len(t, cards, 0)
}

func TestCreateCard(t *testing.T) {
	c := &cards{}
	columnId := uuid.New()
	c.Create("card-1", boardID, columnId)
	c.Create("card-2", boardID, columnId)

	cards, _ := c.List(boardID)
	assert.Len(t, cards, 2)
}

func TestGetCard(t *testing.T) {
	c := &cards{}
	card, _ := c.Create("card-1", boardID, uuid.New())

	_, err := c.Get(uuid.New())
	assert.Error(t, err)

	_, err = c.Get(card.ID)
	assert.NoError(t, err)
}

func TestUpdateCard(t *testing.T) {
	c := &cards{}
	card, _ := c.Create("card-1", boardID, uuid.New())

	card.Name = "card-1-updated"
	err := c.Update(card)
	assert.NoError(t, err)

	ca, _ := c.Get(card.ID)
	assert.Equal(t, "card-1-updated", ca.Name)
}

func TestDeleteCard(t *testing.T) {
	c := &cards{}
	columnId := uuid.New()

	nc, _ := c.Create("card-1", boardID, columnId)
	c.Create("card-2", boardID, columnId)

	err := c.Delete(nc.ID)
	assert.NoError(t, err)

	cards, _ := c.List(boardID)
	assert.Len(t, cards, 1)
}
