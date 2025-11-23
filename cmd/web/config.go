package main

import (
	"flag"
	"fmt"
	"os"
)

// config stores all configurable values
type config struct {
	port           int
	staticDir      string
	secret         string
	initialColumns string
	natsUrl        string
	natsCreds      string
	secure         bool
	enableTimer    bool
	enableStandup  bool
}

func parseConfig() config {
	conf := config{}
	flag.IntVar(&conf.port, "port", 8080, "Port to listen")
	flag.StringVar(&conf.staticDir, "staticDir", "./web/public", "Directory of static files")
	flag.StringVar(&conf.secret, "secret", os.Getenv("GORETRO_SESSION_SECRET"), "Session secret")
	flag.StringVar(&conf.initialColumns, "initialColumns", "Good,Bad,Questions,Emoji", "Initial board columns")
	flag.StringVar(&conf.natsUrl, "nats-url", os.Getenv("GORETRO_NATS_URL"), "NATS Url")
	flag.StringVar(&conf.natsCreds, "nats-cred", os.Getenv("GORETRO_NATS_CREDS"), "Based64 encoded NATS Credentials")
	flag.BoolVar(&conf.secure, "secure", false, "Secure cookie by default")
	flag.BoolVar(&conf.enableTimer, "timer", true, "Enable timer feature")
	flag.BoolVar(&conf.enableStandup, "standup", true, "Enable standup feature")
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
