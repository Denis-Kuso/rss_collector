package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Denis-Kuso/rss_aggregator_p/internal/database"
)

func (s *stateConfig) GetAllFollowedFeeds(w http.ResponseWriter, r *http.Request, user database.User) {

	var errMsg string
	feedFollows, err := s.DB.GetFeedFollowsForUser(r.Context(), user.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errMsg = fmt.Sprintf("no followed feeds found")
			respondWithJSON(w, http.StatusOK, errMsg)
			return
		}
		errMsg = "err retrieving all feedFollows"
		log.Printf("%s: %v\n", errMsg, err)
		respondWithError(w, http.StatusInternalServerError, "can't retrieve feedfollows")
		return
	}
	respondWithJSON(w, http.StatusOK, feedFollows)
	return
}
