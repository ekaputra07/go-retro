package storage

// import (
// 	"sync"

// 	"github.com/ekaputra07/go-retro/internal/model"
// 	"github.com/google/uuid"
// )

// type MemoryStore struct {
// 	boards sync.Map
// }

// func (m *MemoryStore) GetOrCreateBoard(id string) (*model.Board, error) {
// 	newBoard := model.NewBoard(id)
// 	b, _ := m.boards.LoadOrStore(id, newBoard)
// 	return b.(*model.Board), nil
// }

// func (m *MemoryStore) AddColumn(boardID string, column *model.Column) error {
// 	b, err := m.GetOrCreateBoard(boardID)
// 	if err != nil {
// 		return err
// 	}
// 	b.Columns = append(b.Columns, column)
// 	return nil
// }

// func (m *MemoryStore) DeleteColumn(boardID string, columID uuid.UUID) error {
// 	b, err := m.GetOrCreateBoard(boardID)
// 	if err != nil {
// 		return err
// 	}
// 	columns := []*model.Column{}
// 	for _, c := range b.Columns {
// 		if c.ID != columID {
// 			columns = append(columns, c)
// 		}
// 	}
// 	b.Columns = columns
// 	return nil
// }

// func (m *MemoryStore) UpdateColumn(boardID string, column *model.Column) error {
// 	b, err := m.GetOrCreateBoard(boardID)
// 	if err != nil {
// 		return err
// 	}
// 	for _, c := range b.Columns {
// 		if c.ID == column.ID {
// 			c = column
// 			break
// 		}
// 	}
// 	return nil
// }
