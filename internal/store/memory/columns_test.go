package memory

import (
	"context"
	"testing"

	"github.com/ekaputra07/go-retro/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestListColumn(t *testing.T) {
	c := &columns{}
	columns, _ := c.List(context.Background())
	assert.Len(t, columns, 0)
}

func TestCreateColumn(t *testing.T) {
	ctx := context.Background()

	c := &columns{}
	c.Create(ctx, models.NewColumn("col1", 1, boardID))
	c.Create(ctx, models.NewColumn("col2", 2, boardID))

	columns, _ := c.List(ctx)
	assert.Len(t, columns, 2)
}

func TestGetColumn(t *testing.T) {
	ctx := context.Background()

	c := &columns{}
	col := models.NewColumn("col1", 1, boardID)
	c.Create(ctx, col)

	_, err := c.Get(ctx, uuid.New())
	assert.Error(t, err)

	_, err = c.Get(ctx, col.ID)
	assert.NoError(t, err)
}

func TestUpdateColumn(t *testing.T) {
	ctx := context.Background()

	c := &columns{}
	col := models.NewColumn("col1", 1, boardID)
	err := c.Create(ctx, col)
	assert.NoError(t, err)

	col.Name = "col-1-updated"
	err = c.Update(ctx, col)
	assert.NoError(t, err)

	cl, _ := c.Get(ctx, col.ID)
	assert.Equal(t, "col-1-updated", cl.Name)
}

func TestDeleteColumn(t *testing.T) {
	ctx := context.Background()

	c := &columns{}
	col := models.NewColumn("col1", 1, boardID)
	c.Create(ctx, col)
	c.Create(ctx, models.NewColumn("col2", 2, boardID))

	err := c.Delete(ctx, col.ID)
	assert.NoError(t, err)

	columns, _ := c.List(ctx)
	assert.Len(t, columns, 1)
}
