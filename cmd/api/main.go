package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"log"
	"time"

	_ "github.com/lib/pq"
	"url.shortener/internal/data"
	jsonlog "url.shortener/internal/jsonlog"
)

type storageConfig struct {
	storage_type string
	db           databaseConfig
}

type databaseConfig struct {
	dsn          string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type config struct {
	ip           string
	port         int
	baseUrl      string
	suffixLength int
	limiter      struct {
		rps     float64
		burst   int
		enabled bool
	}
	storage storageConfig
}

type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
}

func main() {
	cfg := parseFlags()

	logger := jsonlog.New(log.Writer(), jsonlog.LevelInfo)

	if err := validateConfig(cfg); err != nil {
		logger.PrintFatal(err, nil)
	}

	var models *data.Models
	if cfg.storage.storage_type == "in-memory" {
		models = data.NewModelsInMemory(cfg.baseUrl, cfg.suffixLength)
	} else {
		db, err := openDB(cfg)
		if err != nil {
			logger.PrintFatal(err, nil)
		}
		defer db.Close()

		logger.PrintInfo("Database connection pool established", nil)

		models = data.NewModelsPostgres(cfg.baseUrl, cfg.suffixLength, db)
	}

	app := &application{
		config: cfg,
		logger: logger,
		models: *models,
	}

	err := app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}

func parseFlags() config {
	var cfg config

	flag.StringVar(&cfg.ip, "ip", "0.0.0.0", "The server IP address")
	flag.IntVar(&cfg.port, "port", 50051, "The server port")

	flag.StringVar(&cfg.baseUrl, "base_url", "", "The base URL for short links")
	flag.IntVar(&cfg.suffixLength, "suffix_length", 10, "The length of the URL suffix")

	flag.StringVar(&cfg.storage.storage_type, "storage_type", "in-memory", "The storage type to use for generated URLs (in-memory|postgres)")

	flag.StringVar(&cfg.storage.db.dsn, "db-dsn", "", "PostgreSQL DSN")

	flag.IntVar(&cfg.storage.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.storage.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.storage.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", false, "Enable rate limiter")

	flag.Parse()

	return cfg
}

func validateConfig(cfg config) error {
	if cfg.baseUrl == "" {
		return errors.New("base URL is required")
	}

	if cfg.storage.storage_type != "in-memory" && cfg.storage.storage_type != "postgres" {
		return errors.New("invalid storage type")
	}

	if cfg.storage.storage_type == "postgres" && cfg.storage.db.dsn == "" {
		return errors.New("DSN is required for PostgreSQL storage")
	}

	return nil
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.storage.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.storage.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.storage.db.maxIdleConns)

	duration, err := time.ParseDuration(cfg.storage.db.maxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
