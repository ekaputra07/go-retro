package memstore

import (
	"context"
	"fmt"
	"sync"

	"github.com/ekaputra07/go-retro/internal/models"
	"github.com/google/uuid"
)

// users stores all users in memory
type users struct {
	sync.Map
}

func (u *users) Create(_ context.Context, user models.User) error {
	u.Store(user.ID, user)
	return nil
}

func (u *users) Get(_ context.Context, id uuid.UUID) (*models.User, error) {
	value, ok := u.Load(id)
	if !ok {
		return nil, fmt.Errorf("user id=%s does not exist", id)
	}

	user := value.(models.User)
	return &user, nil
}

func (u *users) Update(_ context.Context, user models.User) error {
	if _, loaded := u.LoadOrStore(user.ID, user); !loaded {
		return fmt.Errorf("user id=%s does not exist", user.ID)
	}
	return nil
}
