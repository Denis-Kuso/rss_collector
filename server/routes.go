package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

// setupGeneralRoutes sets up routes that are not part of the API (health check, root, etc.)
func (a *app) setupRoutes() *chi.Mux {
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

	r.Mount("/v1", a.apiRouter())

	return r
}

// apiRouter sets up the API routes
func (a *app) apiRouter() chi.Router {
	apiRouter := chi.NewRouter()

	apiRouter.Post("/users", a.CreateUser)
	apiRouter.Get("/users", a.MiddlewareAuth(a.GetUserData))
	apiRouter.Get("/feeds", a.GetFeeds)
	apiRouter.Post("/feeds", a.MiddlewareAuth(a.CreateFeed))
	apiRouter.Post("/feed_follows", a.MiddlewareAuth(a.FollowFeed))
	apiRouter.Delete("/feed_follows/{feedFollowID}", a.MiddlewareAuth(a.UnfollowFeed))
	apiRouter.Get("/feed_follows", a.MiddlewareAuth(a.GetAllFollowedFeeds))
	apiRouter.Get("/posts", a.MiddlewareAuth(a.GetPostsFromUser))

	return apiRouter
}
