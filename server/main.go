package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Denis-Kuso/rss_collector/server/internal/database"
	_ "github.com/lib/pq" // importing for side effects
)

func main() {
	// TODO will modify/enrich
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	c := newConfig()
	if showVersion {
		fmt.Printf("Version:\t%s\n", version)
		os.Exit(0)
	}
	db, err := openDB(c)
	defer db.Close()
	if err != nil {
		logger.Fatalf("cannot start db pool: %v", err)
	}
	dQueries := database.New(db)
	a := app{
		cfg:    c,
		logger: logger,
		db:     dQueries,
	}
	err = a.serve()
	if err != nil {
		logger.Fatal(err)
	}
}
