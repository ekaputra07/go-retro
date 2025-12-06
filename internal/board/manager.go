package board

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ekaputra07/go-retro/internal/models"
	"github.com/ekaputra07/go-retro/internal/natsutil"
	"github.com/ekaputra07/go-retro/internal/store"
	"github.com/google/uuid"
)

// BoardManager provides apis to work with board and timer instances.
type BoardManager struct {
	logger              *slog.Logger
	store               *store.Store
	nats                *natsutil.NATS
	timers              map[*timer]bool
	initialBoardColumns []string
	stopped             bool
}

// Healthy returns whether the manager is still running
func (m *BoardManager) Healthy() bool {
	return !m.stopped
}

// Start starts the board manager
func (m *BoardManager) Start(ctx context.Context) {
	m.logger.Info("board-manager running...")

	<-ctx.Done()
	m.logger.Info("Stopping all timers:")
	for t := range m.timers {
		m.logger.Info(fmt.Sprintf("stopping timer %s...", t.BoardID))
		close(t.stopChan)
	}
	m.stopped = true
	m.logger.Info("board-manager stopped")
}

// StartTimer starts the timer process for given boardID, returns true if new process started.
func (m *BoardManager) StartTimer(boardID uuid.UUID) bool {
	_, err := queryTimerStatus(m.nats.Conn, boardID)

	// if target timer not running anywhere, start new one
	if err != nil {
		timer := newTimer(boardID, m.nats, m.logger)

		go timer.run()
		m.timers[timer] = true
		return true
	}
	return false
}

// GetOrCreateBoard get or creates board record
func (m *BoardManager) GetOrCreateBoard(ctx context.Context, id uuid.UUID) (*models.Board, error) {
	b, err := m.store.Boards.Get(ctx, id)

	if b != nil && err == nil {
		// exist
		m.logger.Info("board record exist", "id", id)
		return b, nil
	} else {
		// not found? create new board record with their initial columns
		nb := models.NewBoard(id)
		err = m.store.Boards.Create(ctx, nb)
		if err != nil {
			return nil, err
		}
		m.logger.Info("board record created", "id", id)

		// create initial columns
		for i, c := range m.initialBoardColumns {
			col := models.NewColumn(c, id)
			col.CreatedAt += int64(i) // alter created_at to keep order
			err = m.store.Columns.Create(ctx, col)
			if err != nil {
				return nil, err
			}
			m.logger.Info("board colum created", "name", c)
		}
		return &nb, nil
	}
}

// NewBoardManager creates a new board manager instance
func NewBoardManager(logger *slog.Logger, nats_ *natsutil.NATS, store *store.Store, initialcolumns []string) *BoardManager {
	return &BoardManager{
		logger:              logger,
		nats:                nats_,
		store:               store,
		timers:              make(map[*timer]bool),
		initialBoardColumns: initialcolumns,
		stopped:             false,
	}
}
