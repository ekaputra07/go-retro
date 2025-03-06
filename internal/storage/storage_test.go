package storage

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var boardID = uuid.New()

func testStorageImpl(t *testing.T, init func() Storage) {
	// user
	testCreateUser(t, init)
	testGetUser(t, init)
	testUpdateUser(t, init)
	// board
	testListBoard(t, init)
	testCreateBoard(t, init)
	testGetBoard(t, init)
	testDeleteBoard(t, init)
	// column
	testListColumn(t, init)
	testCreateColumn(t, init)
	testGetColumn(t, init)
	testUpdateColumn(t, init)
	testDeleteColumn(t, init)
	// card
	testListCard(t, init)
	testCreateCard(t, init)
	testGetCard(t, init)
	testUpdateCard(t, init)
	testDeleteCard(t, init)
}

func testCreateUser(t *testing.T, init func() Storage) {
	s := init()
	_, err := s.CreateUser()
	assert.NoError(t, err)
}

func testGetUser(t *testing.T, init func() Storage) {
	s := init()
	u, err := s.CreateUser()

	_, err = s.GetUser(u.ID)
	assert.NoError(t, err)
}

func testUpdateUser(t *testing.T, init func() Storage) {
	s := init()
	u, err := s.CreateUser()
	u.Name = "New Name"

	err = s.UpdateUser(u)
	assert.NoError(t, err)
}

func testListBoard(t *testing.T, init func() Storage) {
	s := init()
	boards, _ := s.ListBoard()
	assert.Len(t, boards, 0)
}

func testCreateBoard(t *testing.T, init func() Storage) {
	s := init()

	s.CreateBoard(boardID)
	boards, _ := s.ListBoard()
	assert.Len(t, boards, 1)
}

func testGetBoard(t *testing.T, init func() Storage) {
	s := init()
	b, _ := s.CreateBoard(boardID)

	_, err := s.GetBoard(uuid.New())
	assert.Error(t, err)

	_, err = s.GetBoard(b.ID)
	assert.NoError(t, err)
}

func testDeleteBoard(t *testing.T, init func() Storage) {
	s := init()
	s.CreateBoard(boardID)
	s.DeleteBoard(boardID)

	boards, _ := s.ListBoard()
	assert.Len(t, boards, 0)
}

func testListColumn(t *testing.T, init func() Storage) {
	s := init()
	columns, _ := s.ListColumn(boardID)
	assert.Len(t, columns, 0)
}

func testCreateColumn(t *testing.T, init func() Storage) {
	s := init()
	s.CreateColumn("col-1", boardID)
	s.CreateColumn("col-2", boardID)

	columns, _ := s.ListColumn(boardID)
	assert.Len(t, columns, 2)
}

func testGetColumn(t *testing.T, init func() Storage) {
	s := init()
	col, _ := s.CreateColumn("col-1", boardID)

	_, err := s.GetColumn(uuid.New())
	assert.Error(t, err)

	_, err = s.GetColumn(col.ID)
	assert.NoError(t, err)
}

func testUpdateColumn(t *testing.T, init func() Storage) {
	s := init()
	col, _ := s.CreateColumn("col-1", boardID)

	col.Name = "col-1-updated"
	err := s.UpdateColumn(col)
	assert.NoError(t, err)

	c, _ := s.GetColumn(col.ID)
	assert.Equal(t, "col-1-updated", c.Name)
}

func testDeleteColumn(t *testing.T, init func() Storage) {
	s := init()
	col, _ := s.CreateColumn("col-1", boardID)
	_, _ = s.CreateColumn("col-2", boardID)

	err := s.DeleteColumn(col.ID)
	assert.NoError(t, err)

	columns, _ := s.ListColumn(boardID)
	assert.Len(t, columns, 1)
}

func testListCard(t *testing.T, init func() Storage) {
	s := init()
	cards, _ := s.ListCard(boardID)
	assert.Len(t, cards, 0)
}

func testCreateCard(t *testing.T, init func() Storage) {
	s := init()
	columnId := uuid.New()
	s.CreateCard("card-1", boardID, columnId)
	s.CreateCard("card-2", boardID, columnId)

	cards, _ := s.ListCard(boardID)
	assert.Len(t, cards, 2)
}

func testGetCard(t *testing.T, init func() Storage) {
	s := init()
	card, err := s.CreateCard("card-1", boardID, uuid.New())

	_, err = s.GetCard(uuid.New())
	assert.Error(t, err)

	_, err = s.GetCard(card.ID)
	assert.NoError(t, err)
}

func testUpdateCard(t *testing.T, init func() Storage) {
	s := init()
	card, _ := s.CreateCard("card-1", boardID, uuid.New())

	card.Name = "card-1-updated"
	err := s.UpdateCard(card)
	assert.NoError(t, err)

	c, _ := s.GetCard(card.ID)
	assert.Equal(t, "card-1-updated", c.Name)
}

func testDeleteCard(t *testing.T, init func() Storage) {
	s := init()
	columnId := uuid.New()
	c, _ := s.CreateCard("card-1", boardID, columnId)
	s.CreateCard("card-2", boardID, columnId)

	err := s.DeleteCard(c.ID)
	assert.NoError(t, err)

	cards, _ := s.ListCard(boardID)
	assert.Len(t, cards, 1)
}
