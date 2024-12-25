package storage

import "testing"

func TestMemoryStore(t *testing.T) {
	testStorageImpl(t, func() Storage {
		return NewMemoryStore()
	})
}
