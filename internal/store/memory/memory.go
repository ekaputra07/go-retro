package memory

import "github.com/ekaputra07/go-retro/internal/store"

func NewMemoryStore() *store.Store {
	return &store.Store{
		Users:   &users{},
		Boards:  &boards{},
		Columns: &columns{},
		Cards:   &cards{},
	}
}
