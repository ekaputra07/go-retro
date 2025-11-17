package memory

import (
	"github.com/ekaputra07/go-retro/internal/store"
)

// NOTE: Memory store is unused in GoRetro V2
// this package left here just for historical reason

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
