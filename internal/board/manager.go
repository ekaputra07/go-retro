package board

import (
	"context"
	"log"

	"github.com/ekaputra07/go-retro/internal/storage"
	"github.com/google/uuid"
)

// BoardManager manages board instances
type BoardManager struct {
	db              storage.Storage
	boards          map[*Board]bool
	registerBoard   chan *Board
	unregisterBoard chan *Board
}

func (bm *BoardManager) Start(ctx context.Context) {
	log.Println("board-manager started")
	for {
		select {
		case b := <-bm.registerBoard:
			bm.boards[b] = true
			log.Printf("board=%s registered", b.ID)
		case b := <-bm.unregisterBoard:
			delete(bm.boards, b)
			log.Printf("board=%s unregistered", b.ID)
		case <-ctx.Done():
			log.Println("board-manager stopped")
			return
		}
	}
}

func (bm *BoardManager) GetOrStartBoard(id uuid.UUID) *Board {
	// if board is running, return it
	for b := range bm.boards {
		if b.ID == id {
			log.Printf("board=%s still running\n", b.ID)
			return b
		}
	}

	board, _ := getOrCreateBoard(id, bm)
	bm.registerBoard <- board
	go board.start()
	return board
}

func NewBoardManager(db storage.Storage) *BoardManager {
	return &BoardManager{
		db:              db,
		boards:          make(map[*Board]bool),
		registerBoard:   make(chan *Board),
		unregisterBoard: make(chan *Board),
	}
}
