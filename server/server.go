package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (a *app) serve() error {
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", a.cfg.port),
		ReadHeaderTimeout: 500 * time.Millisecond,
		ReadTimeout:       500 * time.Millisecond,
		IdleTimeout:       1000 * time.Millisecond,
		Handler:           a.setupRoutes(),
		ErrorLog:          log.Default(), // TODO replace with custom logger
	}
	done := make(chan struct{})
	go worker(done, a.feeds, time.Duration(a.cfg.fetch.reqInterval)*time.Second)
	shutdownErr := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1) // unbuffered chanel might not receive
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit
		close(done)
		slog.Info("server shutdown started", "signal", s.String())
		const gracePeriod time.Duration = 10 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), gracePeriod)
		defer cancel()
		shutdownErr <- server.Shutdown(ctx)
	}()
	slog.Info("starting server", slog.Group("properties", slog.Int("port", a.cfg.port), slog.String("env", a.cfg.env)))
	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	err = <-shutdownErr
	if err != nil { // if shutdown fails
		return err
	}
	slog.Info("server stopped")
	return nil
}
