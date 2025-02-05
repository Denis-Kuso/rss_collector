package main

import (
	"net/http"
)

// setupRoutes creates router for all the routes (public and admin)
func (a *app) setupRoutes() *http.ServeMux {
	r := http.NewServeMux()
	r.Handle("/v1/", http.StripPrefix("/v1", a.apiRouter()))
	// TODO admin router

	return r
}

// apiRouter creates endpoints for the public API.
func (a *app) apiRouter() *http.ServeMux {
	apiRouter := http.NewServeMux()
	apiRouter.HandleFunc("POST /users", a.CreateUser)
	apiRouter.HandleFunc("GET /users", a.MiddlewareAuth(a.GetUserData))
	apiRouter.HandleFunc("GET /feeds", a.GetFeeds)
	apiRouter.HandleFunc("POST /feeds", a.MiddlewareAuth(a.CreateFeed))
	apiRouter.HandleFunc("POST /feed_follows", a.MiddlewareAuth(a.FollowFeed))
	apiRouter.HandleFunc("DELETE /feed_follows/{feedFollowID}", a.MiddlewareAuth(a.UnfollowFeed))
	apiRouter.HandleFunc("GET /feed_follows", a.MiddlewareAuth(a.GetAllFollowedFeeds))
	apiRouter.HandleFunc("GET /posts", a.MiddlewareAuth(a.GetPostsFromUser))

	return apiRouter
}
