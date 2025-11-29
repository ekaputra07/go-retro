// WARNING: NATS store which is based on NATS KV is experimental feature in this project.
// Even though its working, the UX feels noticeably slower due to:
// - NATS KV only support List keys, but the values need to be pulled manually using NATS Get
// - Architectural decision where board status updates always pulls all board columns and cards

package natstore

import (
	"context"
	"time"

	"github.com/ekaputra07/go-retro/internal/natsutil"
	"github.com/ekaputra07/go-retro/internal/store"
	"github.com/nats-io/nats.go/jetstream"
)

const TTL = 2 * time.Hour // only available for 2hrs since creation

func getKV(ctx context.Context, nats *natsutil.NATS, namespace string) (jetstream.KeyValue, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return nats.JS.CreateOrUpdateKeyValue(timeoutCtx, jetstream.KeyValueConfig{
		Bucket:   namespace,
		TTL:      TTL,
		MaxBytes: 1024 * 1000 * 100, // 100Mb
	})
}

func NewStore(ctx context.Context, nats *natsutil.NATS, namespace string) (*store.Store, error) {
	kv, err := getKV(ctx, nats, namespace)
	if err != nil {
		return nil, err
	}
	return &store.Store{
		Users:   &users{kv: kv},
		Boards:  &boards{kv: kv},
		Columns: &columns{kv: kv},
		Cards:   &cards{kv: kv},
	}, nil
}
