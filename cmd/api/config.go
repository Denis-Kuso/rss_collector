package main

import (
	"context"
	"database/sql"
	"flag"
	"time"

	"github.com/Denis-Kuso/rss_collector/internal/storage"
)

var (
	version     string
	showVersion bool
)

type config struct {
	port int
	env  string // or int
	db   struct {
		dsn string
	}
	fetch fetchParams
}

type app struct {
	cfg   config
	users storage.UserStore
	feeds storage.FeedStore
	posts storage.PostStore
}

type fetchParams struct {
	numFeeds    uint
	reqInterval uint // interval between repeated requests
}

func newConfig() config {
	var cfg config
	flag.IntVar(&cfg.port, "port", 3000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|production)")
	flag.UintVar(&cfg.fetch.numFeeds, "n_feeds", 3, "number of feeds")
	flag.UintVar(&cfg.fetch.reqInterval, "req", 100, "request interval in seconds")
	// no default DSN val for DB --> supply from ENV or flag
	flag.StringVar(&cfg.db.dsn, "db-dsn", "", "PostgreSQL DSN")
	flag.BoolVar(&showVersion, "version", false, "show version")
	// TODO add config options for DB
	flag.Parse()
	return cfg
}

func openDB(c config) (*sql.DB, error) {
	var db *sql.DB
	const deadline = 5
	db, err := sql.Open("postgres", c.db.dsn)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), deadline*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
