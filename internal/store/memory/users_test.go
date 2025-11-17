package memory

import (
	"context"
	"testing"

	"github.com/ekaputra07/go-retro/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	u := users{}

	user := models.NewUser(1)
	err := u.Create(context.Background(), user)
	assert.NoError(t, err)
}

func TestGetUser(t *testing.T) {
	u := &users{}

	user := models.NewUser(1)
	_ = u.Create(context.Background(), user)
	_, err := u.Get(context.Background(), user.ID)
	assert.NoError(t, err)
}

func TestUpdateUser(t *testing.T) {
	u := &users{}

	user := models.NewUser(1)
	_ = u.Create(context.Background(), user)
	user.Name = "New Name"

	err := u.Update(context.Background(), user)
	assert.NoError(t, err)
}
