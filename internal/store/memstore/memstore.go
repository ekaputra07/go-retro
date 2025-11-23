package memstore

import (
	"github.com/ekaputra07/go-retro/internal/store"
)

func NewGlobalStore() *store.GlobalStore {
	return &store.GlobalStore{
		Users:  &users{},
		Boards: &boards{},
	}
}

func NewBoardStore() *store.BoardStore {
	return &store.BoardStore{
		Columns: &columns{},
		Cards:   &cards{},
	}
}
