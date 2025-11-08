package memory

import (
	"fmt"
	"sync"

	"github.com/ekaputra07/go-retro/internal/models"
	"github.com/google/uuid"
)

// users stores all users in memory
type users struct {
	sync.Map
}

func (u *users) Create(avatarID int) (*models.User, error) {
	nu := models.NewUser(avatarID)
	u.Store(nu.ID, nu)
	return nu, nil
}

func (u *users) Get(id uuid.UUID) (*models.User, error) {
	user, ok := u.Load(id)
	if !ok {
		return nil, fmt.Errorf("user id=%s does not exist", id)
	}
	return user.(*models.User), nil
}

func (u *users) Update(user *models.User) error {
	if _, loaded := u.LoadOrStore(user.ID, user); !loaded {
		return fmt.Errorf("user id=%s does not exist", user.ID)
	}
	return nil
}
