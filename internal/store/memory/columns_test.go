package memory

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestListColumn(t *testing.T) {
	c := &columns{}
	columns, _ := c.List(boardID)
	assert.Len(t, columns, 0)
}

func TestCreateColumn(t *testing.T) {
	c := &columns{}
	c.Create("col-1", boardID)
	c.Create("col-2", boardID)

	columns, _ := c.List(boardID)
	assert.Len(t, columns, 2)
}

func TestGetColumn(t *testing.T) {
	c := &columns{}
	col, _ := c.Create("col-1", boardID)

	_, err := c.Get(uuid.New())
	assert.Error(t, err)

	_, err = c.Get(col.ID)
	assert.NoError(t, err)
}

func TestUpdateColumn(t *testing.T) {
	c := &columns{}
	col, _ := c.Create("col-1", boardID)

	col.Name = "col-1-updated"
	err := c.Update(col)
	assert.NoError(t, err)

	cl, _ := c.Get(col.ID)
	assert.Equal(t, "col-1-updated", cl.Name)
}

func TestDeleteColumn(t *testing.T) {
	c := &columns{}
	col, _ := c.Create("col-1", boardID)
	_, _ = c.Create("col-2", boardID)

	err := c.Delete(col.ID)
	assert.NoError(t, err)

	columns, _ := c.List(boardID)
	assert.Len(t, columns, 1)
}
