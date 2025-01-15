package main

import (
	"log"
	"net/http"
	"time"
)

func (s *StateConfig) serve() error {
	server := &http.Server{
		Addr:              ":" + s.PortNum,
		ReadHeaderTimeout: 500 * time.Millisecond,
		ReadTimeout:       500 * time.Millisecond,
		IdleTimeout:       1000 * time.Millisecond,
		Handler:           s.setupRoutes(),
	}
	go worker(s.DB, s.WorkOpts.WorkersBreak*time.Second, int(s.WorkOpts.NumWorkers))
	log.Printf("Serving on port: %s\n", s.PortNum)
	return server.ListenAndServe()
}
