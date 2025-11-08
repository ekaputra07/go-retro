package board

import (
	"context"
	"log/slog"

	"github.com/ekaputra07/go-retro/internal/store"
	"github.com/google/uuid"
)

// BoardManager manages board instances
type BoardManager struct {
	logger              *slog.Logger
	store               *store.Store
	initialBoardColumns []string
	boards              map[*Board]bool
	registerChan        chan *Board
	unregisterChan      chan *Board
}

// Start starts the board manager goroutine
func (m *BoardManager) Start(ctx context.Context) {
	m.logger.Info("board-manager running...")
	for {
		select {
		case b := <-m.registerChan:
			m.boards[b] = true
			m.logger.Info("board registered", "id", b.ID)
		case b := <-m.unregisterChan:
			delete(m.boards, b)
			m.logger.Info("board unregistered", "id", b.ID)
		case <-ctx.Done():
			m.logger.Info("board-manager stopped.")
			return
		}
	}
}

// RegisterBoard registers a board to the manager
func (m *BoardManager) RegisterBoard(b *Board) {
	m.registerChan <- b
}

// UnregisterBoard unregisters a board from the manager
func (m *BoardManager) UnregisterBoard(b *Board) {
	m.unregisterChan <- b
}

// GetOrStartBoard returns a board instance by ID, if not exist, it will start a new board
func (m *BoardManager) GetOrStartBoard(id uuid.UUID) *Board {
	// if board is running, return it
	for b := range m.boards {
		if b.ID == id {
			m.logger.Info("board still running", "id", b.ID)
			return b
		}
	}

	// if not, create a new board, register it, and start it
	board, _ := getOrCreateBoard(id, m)
	m.RegisterBoard(board)
	board.Start()
	return board
}

// NewBoardManager creates a new board manager instance
func NewBoardManager(logger *slog.Logger, store *store.Store, initialcolumns []string) *BoardManager {
	return &BoardManager{
		logger:              logger,
		store:               store,
		initialBoardColumns: initialcolumns,
		boards:              make(map[*Board]bool),
		registerChan:        make(chan *Board),
		unregisterChan:      make(chan *Board),
	}
}
