package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
)

func (s *stateConfig) GetFeeds(w http.ResponseWriter, r *http.Request) {

	var errMsg string
	feeds, err := s.DB.GetFeeds(r.Context())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errMsg = "no feeds found"
			respondWithJSON(w, http.StatusOK, errMsg)
			return
		}
		errMsg = "could not retrieve feeds"
		log.Printf("%s: %v\n", errMsg, err)
		respondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}
	publicFeeds := dbFeedToPublicFeeds(feeds)
	respondWithJSON(w, http.StatusOK, publicFeeds)
	return
}
