package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	shutdownErr := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1) // unbuffered chanel might not receive
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit
		log.Printf("Shutting down server: %s", s.String())
		const gracePeriod time.Duration = 10 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), gracePeriod)
		defer cancel()
		shutdownErr <- server.Shutdown(ctx)
	}()
	log.Printf("Serving on port: %s\n", s.PortNum)
	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	err = <-shutdownErr
	if err != nil { // if shutdown fails
		return err
	}
	log.Printf("server stopped")
	return nil
}
