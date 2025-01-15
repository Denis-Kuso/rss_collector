package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // importing for side effects
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed loading enviroment: %v", err)
	}
	cfg := NewCfg()

	go worker(cfg.DB, cfg.WorkOpts.WorkersBreak*time.Second, int(cfg.WorkOpts.NumWorkers))

	server := &http.Server{
		Addr:              ":" + cfg.PortNum,
		ReadHeaderTimeout: 500 * time.Millisecond,
		ReadTimeout:       500 * time.Millisecond,
		IdleTimeout:       1000 * time.Millisecond,
		Handler:           cfg.setupRoutes(),
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
