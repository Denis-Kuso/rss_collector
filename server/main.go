package main

import (
	"log"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // importing for side effects
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed loading enviroment: %v", err)
	}
	cfg := NewCfg()
	err = cfg.serve()
}
