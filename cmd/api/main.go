package main

import (
	"errors"
	"flag"
	"log"

	"url.shortener/internal/data"
	jsonlog "url.shortener/internal/jsonlog"
)

type config struct {
	ip      string
	port    int
	baseUrl string
	storage struct {
		storage_type string
		dsn          string
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

	flag.StringVar(&cfg.baseUrl, "base_url", "", "The base URL for short links")

	flag.StringVar(&cfg.storage.storage_type, "storage_type", "in-memory", "The storage type to use for generated URLs (in-memory, postgres)")
	flag.StringVar(&cfg.storage.dsn, "dsn", "", "PostgreSQL DSN")

	flag.Parse()

	logger := jsonlog.New(log.Writer(), jsonlog.LevelInfo)

	if cfg.baseUrl == "" {
		logger.PrintFatal(errors.New("base URL is required"), nil)
	}

	if cfg.storage.storage_type != "in-memory" && cfg.storage.storage_type != "postgres" {
		logger.PrintFatal(errors.New("invalid storage type"), nil)
	}

	if cfg.storage.storage_type == "postgres" && cfg.storage.dsn == "" {
		logger.PrintFatal(errors.New("DSN is required for PostgreSQL storage"), nil)
	}

	app := &application{
		config: cfg,
		logger: logger,
		models: *data.NewModelsInMemory(cfg.baseUrl),
	}

	err := app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}
