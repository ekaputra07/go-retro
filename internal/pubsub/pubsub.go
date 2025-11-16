package pubsub

// PubSub is an interface for publish/subscribe system
type PubSub interface {
	Subscribe(topic string) chan any
	Unsubscribe(topic string, ch chan any)
	Publish(topic string, message any)
	Close()
}
