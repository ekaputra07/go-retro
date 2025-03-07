package board

import (
	"context"
	"log"

	"github.com/ekaputra07/go-retro/internal/storage"
	"github.com/google/uuid"
)

// BoardManager manages board instances
type BoardManager struct {
	db             storage.Storage
	boards         map[*Board]bool
	registerChan   chan *Board
	unregisterChan chan *Board
}

// Start starts the board manager goroutine
func (bm *BoardManager) Start(ctx context.Context) {
	log.Println("board-manager running...")
	for {
		select {
		case b := <-bm.registerChan:
			bm.boards[b] = true
			log.Printf("board=%s registered", b.ID)
		case b := <-bm.unregisterChan:
			delete(bm.boards, b)
			log.Printf("board=%s unregistered", b.ID)
		case <-ctx.Done():
			log.Println("board-manager stopped.")
			return
		}
	}
}

// RegisterBoard registers a board to the manager
func (bm *BoardManager) RegisterBoard(b *Board) {
	bm.registerChan <- b
}

// UnregisterBoard unregisters a board from the manager
func (bm *BoardManager) UnregisterBoard(b *Board) {
	bm.unregisterChan <- b
}

// GetOrStartBoard returns a board instance by ID, if not exist, it will start a new board
func (bm *BoardManager) GetOrStartBoard(id uuid.UUID) *Board {
	// if board is running, return it
	for b := range bm.boards {
		if b.ID == id {
			log.Printf("board=%s still running\n", b.ID)
			return b
		}
	}

	// if not, create a new board, register it, and start it
	board, _ := getOrCreateBoard(id, bm)
	bm.RegisterBoard(board)
	board.Start()
	return board
}

// NewBoardManager creates a new board manager instance
func NewBoardManager(db storage.Storage) *BoardManager {
	return &BoardManager{
		db:             db,
		boards:         make(map[*Board]bool),
		registerChan:   make(chan *Board),
		unregisterChan: make(chan *Board),
	}
}
