package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	//"github.com/go-chi/cors"
	//"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)
func main() {
	const PORT string = "PORT"
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed loading enviroment.")
	}
	port := os.Getenv(PORT)
	r := chi.NewRouter()
	apiRouter := chi.NewRouter()

	r.Mount("/v1", apiRouter)
	//r.Use(cors.Handler

	corsMux := MiddlewareCors(r)
	// server
	server := &http.Server{
		Addr: ":" + port,
		Handler: corsMux,
	}
	log.Printf("Serving on port: %s\n", port)
	server.ListenAndServe()
}

func MiddlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		log.Printf("Method: %v; URL:%v", r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
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
