package nats

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ekaputra07/go-retro/internal/models"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go/jetstream"
)

type boards struct {
	kv jetstream.KeyValue
}

func (u *boards) key(id uuid.UUID) string {
	return fmt.Sprintf("boards.%s", id)
}

func (b *boards) List(ctx context.Context) ([]models.Board, error) {
	var boards []models.Board
	lister, err := b.kv.ListKeysFiltered(ctx, "boards.*")
	if err != nil {
		return boards, err
	}
	// TODO: limit results
	for key := range lister.Keys() {
		val, err := b.kv.Get(ctx, key)
		if err != nil {
			continue // skip
		}
		var b models.Board
		if err = json.Unmarshal(val.Value(), &b); err != nil {
			continue // skip
		}
		boards = append(boards, b)
	}
	return boards, nil
}

func (b *boards) Create(ctx context.Context, board models.Board) error {
	key := b.key(board.ID)
	_, err := b.kv.Get(ctx, key)
	if err != nil && !errors.Is(err, jetstream.ErrKeyNotFound) {
		return err
	}
	val, err := json.Marshal(board)
	if err != nil {
		return err
	}
	_, err = b.kv.Put(ctx, key, val)
	return err
}

func (b *boards) Get(ctx context.Context, id uuid.UUID) (*models.Board, error) {
	val, err := b.kv.Get(ctx, b.key(id))
	if err != nil {
		return nil, err
	}
	var board models.Board
	err = json.Unmarshal(val.Value(), &board)
	return &board, err
}

func (b *boards) Delete(ctx context.Context, id uuid.UUID) error {
	return b.kv.Delete(ctx, b.key(id))
}
