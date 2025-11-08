package memory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	u := &users{}

	_, err := u.Create(1)
	assert.NoError(t, err)
}

func TestGetUser(t *testing.T) {
	u := &users{}

	nu, _ := u.Create(1)

	_, err := u.Get(nu.ID)
	assert.NoError(t, err)
}

func TestUpdateUser(t *testing.T) {
	u := &users{}

	nu, _ := u.Create(1)
	nu.Name = "New Name"

	err := u.Update(nu)
	assert.NoError(t, err)
}
