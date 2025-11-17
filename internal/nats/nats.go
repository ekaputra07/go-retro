package nats

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// NATS holds reference to Conn, Jetstream and KeyValue
type NATS struct {
	Conn *nats.Conn
	JS   jetstream.JetStream
}

func (n *NATS) Close() {
	n.Conn.Drain()
}

// Setup setups NATS connection and KV store, it panic when any error occured
func Setup(natsURL string) *NATS {
	nc, err := nats.Connect(natsURL)
	if err != nil {
		panic(err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		panic(err)
	}

	return &NATS{
		Conn: nc,
		JS:   js,
	}
}
