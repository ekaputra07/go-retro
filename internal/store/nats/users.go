package nats

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ekaputra07/go-retro/internal/models"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go/jetstream"
)

type users struct {
	kv jetstream.KeyValue
}

func (u *users) key(id uuid.UUID) string {
	return fmt.Sprintf("users.%s", id)
}

func (u *users) Create(ctx context.Context, user models.User) error {
	b, err := json.Marshal(user)
	if err != nil {
		return err
	}
	_, err = u.kv.Put(ctx, u.key(user.ID), b)
	return err
}

func (u *users) Get(ctx context.Context, id uuid.UUID) (*models.User, error) {
	val, err := u.kv.Get(ctx, u.key(id))
	if err != nil {
		return nil, err
	}

	var user models.User
	err = json.Unmarshal(val.Value(), &user)
	return &user, err
}

func (u *users) Update(ctx context.Context, user models.User) error {
	b, err := json.Marshal(user)
	if err != nil {
		return err
	}

	_, err = u.kv.Put(ctx, u.key(user.ID), b)
	return err
}
