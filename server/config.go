package main

import (
	"database/sql"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Denis-Kuso/rss_collector/server/internal/database"
)

const (
	port                 string = "PORT"
	dbConnection         string = "CONN"
	numOfWorkers         string = "NUM_WORKERS"
	interRequestInterval string = "WORKER_REST"
	feedsToFetch         string = "NUM_FEEDS"
	dbType               string = "postgres"
)

// Maintains DB info, port on which to server and options for
// workers fetching the posts
type StateConfig struct {
	DB       *database.Queries
	PortNum  string
	WorkOpts WorkOptions
}

// Maintains number of workers, time period between fetching and how many
// feeds to fetch at a time. They can be specified in .env. If omitted,
// invalid values provided, they resort to default.
type WorkOptions struct {
	NumWorkers   uint
	WorkersBreak time.Duration
	FeedsToFetch uint
}

func NewCfg() *StateConfig {
	const defaultPort = "8080"
	dbURL, found := os.LookupEnv(dbConnection)
	if !found || dbURL == "" {
		log.Fatalf("database connection not specified: \"%s\"\n", dbURL)
	}
	db, err := sql.Open(dbType, dbURL)
	if err = db.Ping(); err != nil {
		log.Fatalf("db: %s not connected: %v", dbURL, err)
	}
	dbQueries := database.New(db)
	return &StateConfig{
		DB:       dbQueries,
		PortNum:  getOptional(port, defaultPort),
		WorkOpts: *newOptions(),
	}
}
func newOptions() *WorkOptions {
	// default values
	const (
		defNumOfWorkers        uint = 3
		defFeedsToFetch        uint = 3
		defInterRequestSeconds uint = 100
	)
	wopt := WorkOptions{
		NumWorkers:   uint(getOptionalAsInt(numOfWorkers, int(defNumOfWorkers))),
		WorkersBreak: time.Duration(getOptionalAsInt(interRequestInterval, int(defInterRequestSeconds))),
		FeedsToFetch: uint(getOptionalAsInt(feedsToFetch, int(defFeedsToFetch))),
	}
	return &wopt
}
func getOptionalAsInt(name string, defaultVal int) int {
	valueStr := getOptional(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}
func getOptional(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
