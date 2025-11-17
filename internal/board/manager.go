package board

import (
	"context"
	"log/slog"

	"github.com/ekaputra07/go-retro/internal/models"
	"github.com/ekaputra07/go-retro/internal/nats"
	"github.com/ekaputra07/go-retro/internal/store"
	"github.com/google/uuid"
)

// BoardManager manages board instances
type BoardManager struct {
	logger              *slog.Logger
	store               *store.GlobalStore
	nats                *nats.NATS
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
func (m *BoardManager) CreateBoard(ctx context.Context, id uuid.UUID) (*Board, error) {
	var board *Board

	// try to get existing board from DB
	b, err := m.store.Boards.Get(ctx, id)

	if err == nil {
		// exist
		m.logger.Info("board record exist", "id", id)
		board, err = newBoard(ctx, m.nats, m, b)
		if err != nil {
			return nil, err
		}
	} else {
		// not found? create new board record with their initial columns
		// TODO: these store operation should run in transaction (on real database)
		nb := models.NewBoard(id)
		err = m.store.Boards.Create(ctx, nb)
		if err != nil {
			return nil, err
		}
		m.logger.Info("board record created", "id", id)
		// board instance from board model
		board, err = newBoard(ctx, m.nats, m, &nb)
		if err != nil {
			return nil, err
		}
	}

	// if no columns records, create initial columns
	columns, err := board.store.Columns.List(ctx)
	if err != nil {
		return nil, err
	}
	if len(columns) == 0 {
		// create initial columns using in-board store
		for _, c := range m.initialBoardColumns {
			col := models.NewColumn(c, board.store.Columns.NextOrder(), board.ID)
			err = board.store.Columns.Create(ctx, col)
			if err != nil {
				return nil, err
			}
			m.logger.Info("board colum created", "name", c)
		}
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
func (m *BoardManager) GetOrCreateBoardProcess(ctx context.Context, id uuid.UUID) (*Board, error) {
	// board already running, return
	if proc := m.GetBoardProcess(id); proc != nil {
		m.logger.Info("board already running", "id", id)
		return proc, nil
	}

	// not running, create new board instance (and new record if needed)
	b, err := m.CreateBoard(ctx, id)
	if err != nil {
		return nil, err
	}
	b.Start()
	m.RegisterBoard(b)

	return b, nil
}

// NewBoardManager creates a new board manager instance
func NewBoardManager(logger *slog.Logger, nats *nats.NATS, store *store.GlobalStore, initialcolumns []string) *BoardManager {
	return &BoardManager{
		logger:              logger,
		nats:                nats,
		store:               store,
		initialBoardColumns: initialcolumns,
		boards:              make(map[*Board]bool),
		registerChan:        make(chan *Board),
		unregisterChan:      make(chan *Board),
		stopped:             false,
	}
}
