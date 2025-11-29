package board

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/ekaputra07/go-retro/internal/models"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
)

func Test_newStream_columns(t *testing.T) {
	t.Run("delete op", func(t *testing.T) {
		s, err := newStream("boards.b.columns.c", jetstream.KeyValueDelete, nil)
		assert.NoError(t, err)
		assert.Equal(t, s.ID, "c")
		assert.Equal(t, s.Op, "del")
		assert.Equal(t, s.Type, "columns")
		assert.Equal(t, s.Object, nil)
	})

	t.Run("purge op", func(t *testing.T) {
		s, err := newStream("boards.b.columns.c", jetstream.KeyValuePurge, nil)
		assert.NoError(t, err)
		assert.Equal(t, s.ID, "c")
		assert.Equal(t, s.Op, "del")
		assert.Equal(t, s.Type, "columns")
		assert.Equal(t, s.Object, nil)
	})

	t.Run("put op", func(t *testing.T) {
		col := models.NewColumn("test", uuid.New())
		val, _ := json.Marshal(col)

		s, err := newStream("boards.b.columns.c", jetstream.KeyValuePut, val)
		assert.NoError(t, err)
		assert.Equal(t, s.ID, "c")
		assert.Equal(t, s.Op, "put")
		assert.Equal(t, s.Type, "columns")
		assert.Equal(t, s.Object, col)
	})
}

func Test_newStream_cards(t *testing.T) {
	t.Run("delete op", func(t *testing.T) {
		s, err := newStream("boards.b.cards.c", jetstream.KeyValueDelete, nil)
		assert.NoError(t, err)
		assert.Equal(t, s.ID, "c")
		assert.Equal(t, s.Op, "del")
		assert.Equal(t, s.Type, "cards")
		assert.Equal(t, s.Object, nil)
	})

	t.Run("purge op", func(t *testing.T) {
		s, err := newStream("boards.b.cards.c", jetstream.KeyValuePurge, nil)
		assert.NoError(t, err)
		assert.Equal(t, s.ID, "c")
		assert.Equal(t, s.Op, "del")
		assert.Equal(t, s.Type, "cards")
		assert.Equal(t, s.Object, nil)
	})

	t.Run("put op", func(t *testing.T) {
		card := models.NewCard("test", uuid.New(), uuid.New())
		val, _ := json.Marshal(card)

		s, err := newStream("boards.b.cards.c", jetstream.KeyValuePut, val)
		assert.NoError(t, err)
		assert.Equal(t, s.ID, "c")
		assert.Equal(t, s.Op, "put")
		assert.Equal(t, s.Type, "cards")
		assert.Equal(t, s.Object, card)
	})
}

func Test_newStream_others(t *testing.T) {
	t.Run("delete op", func(t *testing.T) {
		s, err := newStream("boards.b.unknown.c", jetstream.KeyValueDelete, nil)
		assert.NoError(t, err)
		assert.Equal(t, s.ID, "c")
		assert.Equal(t, s.Op, "del")
		assert.Equal(t, s.Type, "unknown")
		assert.Equal(t, s.Object, nil)
	})

	t.Run("purge op", func(t *testing.T) {
		s, err := newStream("boards.b.unknown.c", jetstream.KeyValuePurge, nil)
		assert.NoError(t, err)
		assert.Equal(t, s.ID, "c")
		assert.Equal(t, s.Op, "del")
		assert.Equal(t, s.Type, "unknown")
		assert.Equal(t, s.Object, nil)
	})

	t.Run("put op", func(t *testing.T) {
		s, err := newStream("boards.b.unknown.c", jetstream.KeyValuePut, nil)
		assert.Nil(t, s)
		assert.Error(t, err)
	})
}

func Test_message_dataGet(t *testing.T) {
	t.Run("infer error", func(t *testing.T) {
		m := message{messageTypeColumnNew, nil, models.NewUser(1)}
		data, err := m.dataGet("id")
		assert.Nil(t, data)
		assert.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), "can't convert"))
	})

	t.Run("key error", func(t *testing.T) {
		d := map[string]any{
			"id": "test",
		}
		m := message{messageTypeColumnNew, d, models.NewUser(1)}
		data, err := m.dataGet("name")
		assert.Nil(t, data)
		assert.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), "no key"))
	})

	t.Run("ok", func(t *testing.T) {
		d := map[string]any{
			"id": "test",
		}
		m := message{messageTypeColumnNew, d, models.NewUser(1)}
		data, err := m.dataGet("id")
		assert.NoError(t, err)
		assert.Equal(t, "test", data)
	})
}

func Test_message_stringVar(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		d := map[string]any{
			"id": "test",
		}
		var val string
		m := message{messageTypeColumnNew, d, models.NewUser(1)}
		err := m.stringVar(&val, "id")
		assert.NoError(t, err)
		assert.Equal(t, "test", val)
	})

	t.Run("infer error", func(t *testing.T) {
		d := map[string]any{
			"id": 123,
		}
		var val string
		m := message{messageTypeColumnNew, d, models.NewUser(1)}
		err := m.stringVar(&val, "id")
		assert.Error(t, err)
		assert.Equal(t, "", val)
	})
}

func Test_message_intVar(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		d := map[string]any{
			"id": float64(123),
		}
		var val int
		m := message{messageTypeColumnNew, d, models.NewUser(1)}
		err := m.intVar(&val, "id")
		assert.NoError(t, err)
		assert.Equal(t, 123, val)
	})

	t.Run("infer error", func(t *testing.T) {
		d := map[string]any{
			"id": "123",
		}
		var val int
		m := message{messageTypeColumnNew, d, models.NewUser(1)}
		err := m.intVar(&val, "id")
		assert.Error(t, err)
		assert.Equal(t, 0, val)
	})
}

func Test_message_uuidVar(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		id := uuid.New()
		d := map[string]any{
			"id": id.String(),
		}
		var val uuid.UUID
		m := message{messageTypeColumnNew, d, models.NewUser(1)}
		err := m.uuidVar(&val, "id")
		assert.NoError(t, err)
		assert.Equal(t, id, val)
	})

	t.Run("infer error", func(t *testing.T) {
		d := map[string]any{
			"id": "123",
		}
		var val uuid.UUID
		m := message{messageTypeColumnNew, d, models.NewUser(1)}
		err := m.uuidVar(&val, "id")
		assert.Error(t, err)
		assert.Equal(t, uuid.Nil, val)
	})
}
