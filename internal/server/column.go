package server

import "github.com/google/uuid"

var defaultColumns = []string{"Good", "Bad", "Questions"}

type column struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (b *board) columnByID(id uuid.UUID) *column {
	for _, c := range b.Columns {
		if c.ID == id {
			return c
		}
	}
	return nil
}
