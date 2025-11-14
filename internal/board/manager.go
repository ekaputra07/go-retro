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
	stopped             bool
}

// Healthy returns whether the manager is still running
func (m *BoardManager) Healthy() bool {
	return !m.stopped
}

// Start starts the board manager goroutine
func (m *BoardManager) Start(ctx context.Context) {
	m.logger.Info("board-manager running...")

	defer func(m *BoardManager) {
		m.stopped = true
	}(m)

	for {
		select {
		case b := <-m.registerChan:
			m.boards[b] = true
			m.logger.Info("board registered", "id", b.ID)
		case b := <-m.unregisterChan:
			delete(m.boards, b)
			m.logger.Info("board unregistered", "id", b.ID)
		case <-ctx.Done():
			m.logger.Info("board-manager stopped")
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

// CreateBoard creates board instance
func (m *BoardManager) CreateBoard(id uuid.UUID) (*Board, error) {
	// try to get existing board from DB
	b, err := m.store.Boards.Get(id)

	// not found? create new board record with their initial columns
	if err != nil {
		// TODO: these store operation should run in transaction (on real database)
		b, err = m.store.Boards.Create(id)
		if err != nil {
			return nil, err
		}
		for _, c := range m.initialBoardColumns {
			_, err := m.store.Columns.Create(c, b.ID)
			if err != nil {
				return nil, err
			}
		}
	}

	// start and register new board process to the manager
	board := &Board{
		Board:   b,
		manager: m,
		logger:  m.logger,
		store:   m.store,
		clients: make(map[*Client]bool),
		join:    make(chan *Client),
		leave:   make(chan *Client),
		message: make(chan message),
		stop:    make(chan struct{}),
		timer:   newTimer(m.logger),
	}
	return board, nil
}

// GetBoardProcess get running board from manager's boards registry
func (m *BoardManager) GetBoardProcess(id uuid.UUID) *Board {
	for b := range m.boards {
		if b.ID == id {
			return b
		}
	}
	return nil
}

// GetOrCreateBoardProcess returns existing or new board process
func (m *BoardManager) GetOrCreateBoardProcess(id uuid.UUID) (*Board, error) {
	// board already running, return
	if proc := m.GetBoardProcess(id); proc != nil {
		m.logger.Info("board already running", "id", id)
		return proc, nil
	}

	m.logger.Info("new board started", "id", id)
	// not running, create new process (and new record if needed)
	b, err := m.CreateBoard(id)
	if err != nil {
		return nil, err
	}
	b.Start()
	m.RegisterBoard(b)

	return b, nil
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
		stopped:             false,
	}
}
