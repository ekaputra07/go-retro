package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	appName    string = "GoRetro"
	appTagline string = "Minimalist retro board for happy teams 😉"
	appVersion string = os.Getenv("GORETRO_VERSION")
)

// config stores all configurable values
type config struct {
	port           int
	secret         string
	initialColumns string
	natsUrl        string
	natsCreds      string
	secure         bool
}

func parseConfig() config {
	conf := config{}
	flag.IntVar(&conf.port, "port", 8080, "Port to listen")
	flag.StringVar(&conf.secret, "secret", os.Getenv("GORETRO_SESSION_SECRET"), "Session secret")
	flag.StringVar(&conf.initialColumns, "initialColumns", "Good,Bad,Questions,Emoji", "Initial board columns")
	flag.StringVar(&conf.natsUrl, "nats-url", os.Getenv("GORETRO_NATS_URL"), "NATS Url")
	flag.StringVar(&conf.natsCreds, "nats-cred", os.Getenv("GORETRO_NATS_CREDS"), "Based64 encoded NATS Credentials")
	flag.BoolVar(&conf.secure, "secure", false, "Secure cookie by default")
	flag.Parse()

	// make sure secret is not empty
	if conf.secret == "" {
		fmt.Println(
			"Secret is missing!",
			"Set secret via environment variable `GORETRO_SESSION_SECRET` (recommended)",
			"or via `-secret` flag.",
		)
		os.Exit(1)
	}
	return conf
}
