package natsutil

import (
	"encoding/base64"

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

// Connect setups NATS connection and KV store, it panic when any error occured
func Connect(url, credentials string) *NATS {
	options := []nats.Option{
		nats.Name("goretro-web"),
	}

	if credentials != "none" {
		// connection credentials string (base64 encoded)
		credBytes, err := base64.StdEncoding.DecodeString(credentials)
		if err != nil {
			panic(err)
		}
		options = append(options, nats.UserCredentialBytes(credBytes))
	}

	nc, err := nats.Connect(url, options...)
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
