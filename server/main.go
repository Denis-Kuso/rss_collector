package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // importing for side effects
	"github.com/rs/cors"
)

func main() {
	const (
		root              = "/"
		ready             = "/readiness"
		errorEndpoint     = "/err"
		users             = "/users"
		feeds             = "/feeds"
		follow_feeds      = "/feed_follows"
		posts             = "/posts"
		QUERY_FEED_FOLLOW = "feedFollowID"
	)
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed loading enviroment: %v", err)
	}
	cfg := NewCfg()

	r := chi.NewRouter()
	apiRouter := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.URLFormat)
	// Use default options for now
	r.Use(cors.Default().Handler)

	r.Get(root, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome Gandalf"))
	})
	r.Mount("/v1", apiRouter)
	apiRouter.Get(ready, func(w http.ResponseWriter, r *http.Request) {
		respondWithJSON(w, http.StatusOK, "status:ok")
	})
	apiRouter.Get(errorEndpoint, func(w http.ResponseWriter, r *http.Request) {
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
	})

	go worker(cfg.DB, cfg.WorkOpts.WorkersBreak*time.Second, int(cfg.WorkOpts.NumWorkers))

	apiRouter.Post(users, cfg.CreateUser)
	apiRouter.Get(users, cfg.MiddlewareAuth(cfg.GetUserData))
	apiRouter.Get(feeds, cfg.GetFeeds)
	apiRouter.Post(feeds, cfg.MiddlewareAuth(cfg.CreateFeed))
	apiRouter.Post(follow_feeds, cfg.MiddlewareAuth(cfg.FollowFeed))
	apiRouter.Delete(follow_feeds+"/{"+QUERY_FEED_FOLLOW+"}", cfg.MiddlewareAuth(cfg.UnfollowFeed))
	apiRouter.Get(follow_feeds, cfg.MiddlewareAuth(cfg.GetAllFollowedFeeds))
	apiRouter.Get(posts, cfg.MiddlewareAuth(cfg.GetPostsFromUser))
	server := &http.Server{
		Addr:              ":" + cfg.PortNum,
		ReadHeaderTimeout: 500 * time.Millisecond,
		ReadTimeout:       500 * time.Millisecond,
		IdleTimeout:       1000 * time.Millisecond,
		Handler:           r,
	}

	log.Printf("Serving on port: %s\n", cfg.PortNum)
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