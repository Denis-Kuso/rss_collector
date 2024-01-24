package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)
func main() {
	const PORT string = "PORT"
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed loading enviroment.")
	}
	port := os.Getenv(PORT)
	ready := "/readiness"
	errorEndpoint := "/err"

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

	server := &http.Server{
		Addr: ":" + port,
		Handler: r,
	}

	log.Printf("Serving on port: %s\n", port)
	server.ListenAndServe()
	//http.ListenAndServe(":"+port, r)
}

//func MiddlewareCors(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		w.Header().Set("Access-Control-Allow-Origin", "*")
//		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
//		w.Header().Set("Access-Control-Allow-Headers", "*")
//		if r.Method == "OPTIONS" {
//			w.WriteHeader(http.StatusOK)
//			return
//		}
//		log.Printf("Method: %v; URL:%v", r.Method, r.URL)
//		next.ServeHTTP(w, r)
//	})
//}

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
