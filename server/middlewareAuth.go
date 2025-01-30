package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Denis-Kuso/rss_collector/server/internal/auth"
	"github.com/Denis-Kuso/rss_collector/server/internal/storage"
)

type contextKey string

const userIDctx = contextKey("userID")

type authenicatedHandler func(w http.ResponseWriter, r *http.Request)

func (a *app) MiddlewareAuth(handler authenicatedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Check auth header
		APIKey, err := auth.GetAPIKey(r.Header)
		var msg string
		if err != nil {
			if errors.Is(err, auth.ErrNoAuthHeaderIncluded) {
				msg = "no header included"
				respondWithError(w, http.StatusBadRequest, "no header included")
				return
			}
			if errors.Is(err, auth.ErrMalformedAuthHeader) {
				msg = err.Error()
				respondWithError(w, http.StatusBadRequest, msg)
				return
			}
		}
		user, err := a.users.WhoIs(r.Context(), APIKey)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				msg = fmt.Sprintf("you must be authenticated to access this resource")
				respondWithError(w, http.StatusUnauthorized, msg)
				return
			}
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		ctx := context.WithValue(r.Context(), userIDctx, user.ID)
		handler(w, r.WithContext(ctx))
	}
}
