package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Denis-Kuso/rss_collector/server/internal/storage"

	_ "github.com/lib/pq" // importing for side effects
)

func main() {
	logger := slog.NewJSONHandler(os.Stdout, nil)
	noviLoger := slog.New(logger)
	slog.SetDefault(noviLoger)
	c := newConfig()
	if showVersion {
		fmt.Printf("Version:\t%s\n", version)
		os.Exit(0)
	}
	db, err := openDB(c)
	defer db.Close()
	if err != nil {
		slog.Warn("cannot start db pool", "error", err, "dsn", c.db.dsn)
		os.Exit(1)
	}

	m := storage.NewUsersModel(db)
	f := storage.NewFeedsModel(db)
	a := app{
		cfg:   c, // TODO this could be options instead// fetchParams not entire config
		users: m,
		feeds: f,
	}
	err = a.serve()
	if err != nil {
		slog.Warn("server shutdown failure", "error", err)
		os.Exit(1)
	}
}
