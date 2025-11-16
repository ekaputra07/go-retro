package pubsub

import (
	"slices"
	"sync"
)

// Local implementations in-process Pubsub mechanism
type Local struct {
	subscribers map[string][]chan any
	mu          sync.RWMutex
	closed      bool
}

// Subscribe create message channel and sunbscribe to a topic
func (ps *Local) Subscribe(topic string) chan any {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	// buffered channel size of 1 so it doesn't block when sending message to it
	ch := make(chan any, 1)
	ps.subscribers[topic] = append(ps.subscribers[topic], ch)
	return ch
}

// Unsubscribe remove a channel from topic subscribers and close it
func (ps *Local) Unsubscribe(topic string, ch chan any) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	subs := ps.subscribers[topic]
	if index := slices.Index(subs, ch); index != -1 {
		ps.subscribers[topic] = slices.Delete(subs, index, index+1)
		close(ch)
	}
}

// Publish publishes message to a topic
func (ps *Local) Publish(topic string, message any) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	if ps.closed {
		return
	}

	for _, sub := range ps.subscribers[topic] {
		sub <- message
	}
}

// Closes closes all subscribers channel
func (ps *Local) Close() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	for _, subs := range ps.subscribers {
		for _, ch := range subs {
			close(ch)
		}
	}
}

func NewLocalPubsub() *Local {
	return &Local{
		subscribers: make(map[string][]chan any),
		closed:      false,
	}
}
