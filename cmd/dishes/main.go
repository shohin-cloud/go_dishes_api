package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"

	"github.com/shohin-cloud/dishes-api/pkg/dishes/model"
	"github.com/shohin-cloud/dishes-api/pkg/jsonlog"
	"github.com/shohin-cloud/dishes-api/pkg/vcs"

	_ "github.com/lib/pq"
)

var (
	version = vcs.Version()
)

type config struct {
	port string
	env  string
	db   struct {
		dsn string
	}
}

type application struct {
	config config
	models model.Models
	logger *jsonlog.Logger
}

func main() {
	fmt.Println("Started server")
	var cfg config
	flag.StringVar(&cfg.port, "port", ":8060", "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://postgres:rootroot@db:5432/postgres?sslmode=disable", "PostgreSQL DSN")
	flag.Parse()

	// Init logger
	logger := jsonlog.NewLogger(os.Stdout, jsonlog.LevelInfo)

	// Connect to DB
	db, err := openDB(cfg)
	if err != nil {
		logger.PrintError(err, nil)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.PrintFatal(err, nil)
		}
	}()

	app := &application{
		config: cfg,
		models: model.NewModels(db),
		logger: logger,
	}

	if err := app.serve(); err != nil {
		logger.PrintFatal(err, nil)
	}

}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
