package main

import (
	"flag"
	"log"

	jsonlog "url.shortener/internal/jsonlog"
)

type config struct {
	ip   string
	port int
}

type application struct {
	config config
	logger *jsonlog.Logger
}

func main() {
	var cfg config

	flag.StringVar(&cfg.ip, "ip", "localhost", "The server IP address")
	flag.IntVar(&cfg.port, "port", 50051, "The server port")

	flag.Parse()

	logger := jsonlog.New(log.Writer(), jsonlog.LevelInfo)

	app := &application{
		config: cfg,
		logger: logger,
	}

	err := app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}
