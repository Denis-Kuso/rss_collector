package main


import (
	_ "github.com/lib/pq" // importing for side effects
	"encoding/json"
	"log"
	"net/http"
	"os"
	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/Denis-Kuso/rss_aggregator_p/internal/database"
	"database/sql"
//	"fmt"
)

type stateConfig struct {
	DB *database.Queries
}

const (
	QUERY_FEED_FOLLOW = "feedFollowID"
	)

func main() {
	const (
		PORT string = "PORT"
		CONN string = "CONN"
	)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed loading enviroment.")
	}
	port := os.Getenv(PORT)
	dbURL := os.Getenv(CONN)
	if (port == "" || dbURL == "") {
		log.Fatalf("Environment variables undefined\n")
	}
	// Init db
	db, err := sql.Open("postgres", dbURL)
	dbQueries := database.New(db)
	state := stateConfig{dbQueries}

	ready := "/readiness"
	errorEndpoint := "/err"
	users := "/users"
	feeds := "/feeds"
	follow_feeds := "/feed_follows"
	r := chi.NewRouter()
	apiRouter := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.URLFormat)
	// Use default options for now
	r.Use(cors.Default().Handler)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome Gandalf"))
	})
	r.Mount("/v1",apiRouter)
	apiRouter.Get(ready, func(w http.ResponseWriter, r *http.Request) {
		respondWithJSON(w, http.StatusOK, "status:ok" )
		})	
	apiRouter.Get(errorEndpoint, func(w http.ResponseWriter, r *http.Request){
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		})

	apiRouter.Post(users, state.CreateUser)
	apiRouter.Get(users, state.MiddlewareAuth(state.GetUserData))
	apiRouter.Post(feeds, state.MiddlewareAuth(state.CreateFeed))
	apiRouter.Post(follow_feeds, state.MiddlewareAuth(state.FollowFeed))
	apiRouter.Delete(follow_feeds +"/{"+QUERY_FEED_FOLLOW +"}", state.MiddlewareAuth(state.UnfollowFeed))
	apiRouter.Get(follow_feeds, state.MiddlewareAuth(state.GetAllFollowedFeeds))
	server := &http.Server{
		Addr: ":" + port,
		Handler: r,
	}

	log.Printf("Serving on port: %s\n", port)
//	rss, err := URLtoFeed(LANES_BLOG)
//	if err != nil {
//		log.Printf("ERR rss to feed: %v\n", err)
//	}
//	
//	//log.Printf("%v\n",rss)
//	for _, item := range rss.Feed1.Feeds{
//		fmt.Println("----------------------")
//		fmt.Printf("%v: %v\n",item.Title, item.Id)
//	}
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
