package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/Denis-Kuso/rss_collector/server/internal/auth"
)

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
		// TODO still better to getByID
		user, err := a.users.Get(r.Context(), APIKey)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				msg = fmt.Sprintf("no user with APIKey: %s", APIKey)
				respondWithError(w, http.StatusNotFound, msg)
				return
			}
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		// TODO i really want to pass userID
		ctx := context.WithValue(r.Context(), "APIkey", user.APIkey) // TODO ensure type safety
		handler(w, r.WithContext(ctx))
	}
}
