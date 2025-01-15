package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

// setupGeneralRoutes sets up routes that are not part of the API (health check, root, etc.)
func (cfg *StateConfig) setupRoutes() *chi.Mux {
	// Middleware
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.URLFormat)
	r.Use(cors.Default().Handler)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Get("/readiness", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok!"))
	})

	r.Mount("/v1", cfg.apiRouter())

	return r
}

// apiRouter sets up the API routes
func (cfg *StateConfig) apiRouter() chi.Router {
	apiRouter := chi.NewRouter()

	apiRouter.Post("/users", cfg.CreateUser)
	apiRouter.Get("/users", cfg.MiddlewareAuth(cfg.GetUserData))
	apiRouter.Get("/feeds", cfg.GetFeeds)
	apiRouter.Post("/feeds", cfg.MiddlewareAuth(cfg.CreateFeed))
	apiRouter.Post("/feed_follows", cfg.MiddlewareAuth(cfg.FollowFeed))
	apiRouter.Delete("/feed_follows/{feedFollowID}", cfg.MiddlewareAuth(cfg.UnfollowFeed))
	apiRouter.Get("/feed_follows", cfg.MiddlewareAuth(cfg.GetAllFollowedFeeds))
	apiRouter.Get("/posts", cfg.MiddlewareAuth(cfg.GetPostsFromUser))

	return apiRouter
}
