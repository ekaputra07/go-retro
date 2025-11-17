package nats

import (
	"context"
	"time"

	"github.com/ekaputra07/go-retro/internal/nats"
	"github.com/ekaputra07/go-retro/internal/store"
	"github.com/nats-io/nats.go/jetstream"
)

const TTL = 6 * time.Hour // only available for 6hrs since creation

func getKV(ctx context.Context, nats *nats.NATS, namespace string) (jetstream.KeyValue, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return nats.JS.CreateOrUpdateKeyValue(timeoutCtx, jetstream.KeyValueConfig{
		Bucket: namespace,
		TTL:    TTL,
	})

}

func NewGlobalStore(ctx context.Context, nats *nats.NATS, namespace string) (*store.GlobalStore, error) {
	kv, err := getKV(ctx, nats, namespace)
	if err != nil {
		return nil, err
	}
	return &store.GlobalStore{
		Users:  &users{kv: kv},
		Boards: &boards{kv: kv},
	}, nil
}

func NewBoardStore(ctx context.Context, nats *nats.NATS, namespace string) (*store.BoardStore, error) {
	kv, err := getKV(ctx, nats, namespace)
	if err != nil {
		return nil, err
	}
	return &store.BoardStore{
		Columns: &columns{kv: kv},
		Cards:   &cards{kv: kv},
	}, nil
}
