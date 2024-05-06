package main

import (
	"database/sql"
	"encoding/json"
	"github.com/Denis-Kuso/rss_aggregator_p/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // importing for side effects
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
	"time"
)

type stateConfig struct {
	DB *database.Queries
}

const (
)

func main() {
	const (
		PORT string = "PORT"
		CONN string = "CONN"

		ready = "/readiness"
		errorEndpoint = "/err"
		users = "/users"
		feeds = "/feeds"
		follow_feeds = "/feed_follows"
		QUERY_FEED_FOLLOW = "feedFollowID"
		posts = "/posts"
	)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed loading enviroment.")
	}
	port := os.Getenv(PORT)
	dbURL := os.Getenv(CONN)
	if port == "" || dbURL == "" {
		log.Fatalf("Environment variables undefined\n")
	}
	// Init db
	db, err := sql.Open("postgres", dbURL)
	if err := db.Ping(); err != nil {
		log.Fatalf("db not connected: %v", err)
	}
	dbQueries := database.New(db)
	state := stateConfig{dbQueries}

	r := chi.NewRouter()
	apiRouter := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.URLFormat)
	// Use default options for now
	r.Use(cors.Default().Handler)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome Gandalf"))
	})
	r.Mount("/v1", apiRouter)
	apiRouter.Get(ready, func(w http.ResponseWriter, r *http.Request) {
		respondWithJSON(w, http.StatusOK, "status:ok")
	})
	apiRouter.Get(errorEndpoint, func(w http.ResponseWriter, r *http.Request) {
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
	})

	go worker(dbQueries, 10*time.Second, 3)
	apiRouter.Post(users, state.CreateUser)
	apiRouter.Get(users, state.MiddlewareAuth(state.GetUserData))
	apiRouter.Get(feeds, state.GetFeeds)
	apiRouter.Post(feeds, state.MiddlewareAuth(state.CreateFeed))
	apiRouter.Post(follow_feeds, state.MiddlewareAuth(state.FollowFeed))
	apiRouter.Delete(follow_feeds+"/{"+QUERY_FEED_FOLLOW+"}", state.MiddlewareAuth(state.UnfollowFeed))
	apiRouter.Get(follow_feeds, state.MiddlewareAuth(state.GetAllFollowedFeeds))
	apiRouter.Get(posts, state.MiddlewareAuth(state.GetPostsFromUser))
	server := &http.Server{
		Addr:              ":" + port,
		ReadHeaderTimeout: 500 * time.Millisecond,
		ReadTimeout:       500 * time.Millisecond,
		IdleTimeout:       1000 * time.Millisecond,
		Handler:           r,
	}

	log.Printf("Serving on port: %s\n", port)
	server.ListenAndServe()
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}
