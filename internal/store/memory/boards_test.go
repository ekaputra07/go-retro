package memory

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var boardID = uuid.New()

func TestListBoard(t *testing.T) {
	b := &boards{}
	boards, _ := b.List()
	assert.Len(t, boards, 0)
}

func TestCreateBoard(t *testing.T) {
	b := &boards{}

	b.Create(boardID)
	boards, _ := b.List()
	assert.Len(t, boards, 1)
}

func TestGetBoard(t *testing.T) {
	b := &boards{}
	nb, _ := b.Create(boardID)

	_, err := b.Get(uuid.New())
	assert.Error(t, err)

	_, err = b.Get(nb.ID)
	assert.NoError(t, err)
}

func TestDeleteBoard(t *testing.T) {
	b := &boards{}
	nb, _ := b.Create(boardID)
	b.Create(boardID)

	err := b.Delete(nb.ID)
	assert.NoError(t, err)

	boards, _ := b.List()
	assert.Len(t, boards, 0)
}
