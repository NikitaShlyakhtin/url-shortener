package main

import (
	"flag"
	"log"

	"url.shortener/internal/data"
	jsonlog "url.shortener/internal/jsonlog"
)

type config struct {
	ip      string
	port    int
	storage struct {
		inMemory bool
		dsn      string
	}
}

type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
}

func main() {
	var cfg config

	flag.StringVar(&cfg.ip, "ip", "localhost", "The server IP address")
	flag.IntVar(&cfg.port, "port", 50051, "The server port")

	flag.BoolVar(&cfg.storage.inMemory, "in_memory", true, "Use in-memory storage")
	flag.StringVar(&cfg.storage.dsn, "dsn", "", "PostgreSQL DSN")

	flag.Parse()

	logger := jsonlog.New(log.Writer(), jsonlog.LevelInfo)

	app := &application{
		config: cfg,
		logger: logger,
		models: *data.NewModelsInMemory(),
	}

	err := app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}
