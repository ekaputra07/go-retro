package pubsub

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLocalPubsub(t *testing.T) {
	ps := NewLocalPubsub()

	assert.False(t, ps.closed)
	assert.Len(t, ps.subscribers, 0)
}

func TestLocalSubscribe(t *testing.T) {
	ps := NewLocalPubsub()
	sub1 := ps.Subscribe("topic1")
	sub2 := ps.Subscribe("topic2")

	assert.NotEqual(t, sub1, sub2)
	assert.Len(t, ps.subscribers, 2)
	assert.Len(t, ps.subscribers["topic1"], 1)
	assert.Len(t, ps.subscribers["topic2"], 1)
	ps.Close()
}

func TestLocalUnsbscribe(t *testing.T) {
	ps := NewLocalPubsub()
	sub := ps.Subscribe("topic1")
	ps.Unsubscribe("topic0", sub) // do nothing
	ps.Unsubscribe("topic1", sub)
	assert.Len(t, ps.subscribers["topic1"], 0)

	_, ok := <-sub
	assert.False(t, ok)
}

func TestLocalPublish(t *testing.T) {
	ps := NewLocalPubsub()
	sub := ps.Subscribe("topic1")

	ps.Publish("topic1", "hello")
	assert.Equal(t, <-sub, "hello")

	ps.Publish("topic1", "world")
	assert.Equal(t, <-sub, "world")

	ps.Close()
}

func TestLocalClose(t *testing.T) {
	ps := NewLocalPubsub()
	sub := ps.Subscribe("topic")
	sub2 := ps.Subscribe("topic")
	ps.Close()

	_, ok := <-sub
	_, ok2 := <-sub2
	assert.False(t, ok)
	assert.False(t, ok2)
}
