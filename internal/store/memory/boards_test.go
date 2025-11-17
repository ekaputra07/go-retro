package memory

import (
	"context"
	"testing"

	"github.com/ekaputra07/go-retro/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var boardID = uuid.New()

func TestListBoard(t *testing.T) {
	b := &boards{}
	boards, _ := b.List(context.Background())
	assert.Len(t, boards, 0)
}

func TestCreateBoard(t *testing.T) {
	ctx := context.Background()

	b := &boards{}
	nb := models.NewBoard(uuid.New())
	b.Create(ctx, nb)
	boards, _ := b.List(ctx)
	assert.Len(t, boards, 1)
}

func TestGetBoard(t *testing.T) {
	ctx := context.Background()

	b := &boards{}
	nb := models.NewBoard(uuid.New())
	b.Create(ctx, nb)
	_, err := b.Get(ctx, uuid.New())
	assert.Error(t, err)

	_, err = b.Get(ctx, nb.ID)
	assert.NoError(t, err)
}

func TestDeleteBoard(t *testing.T) {
	ctx := context.Background()

	b := &boards{}
	nb := models.NewBoard(uuid.New())
	b.Create(ctx, nb)

	err := b.Delete(ctx, nb.ID)
	assert.NoError(t, err)

	boards, _ := b.List(ctx)
	assert.Len(t, boards, 0)
}
