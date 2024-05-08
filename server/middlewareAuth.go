package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/Denis-Kuso/rss_collector/server/internal/auth"
	"github.com/Denis-Kuso/rss_collector/server/internal/database"
)

type authenicatedHandler func(w http.ResponseWriter, r *http.Request, user database.User)

func (s *StateConfig) MiddlewareAuth(handler authenicatedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Check auth header
		apiKey, err := auth.GetAPIKey(r.Header)
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
		user, err := s.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				msg = fmt.Sprintf("no user with apiKey: %s", apiKey)
				respondWithError(w, http.StatusNotFound, msg)
				return
			}
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		handler(w, r, user)
	}
}
