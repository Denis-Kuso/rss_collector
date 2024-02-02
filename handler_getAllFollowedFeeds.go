package main 

import (
	"log"
	"net/http"
	"github.com/Denis-Kuso/rss_aggregator_p/internal/database"
)

func (s *stateConfig) GetAllFollowedFeeds(w http.ResponseWriter, r *http.Request, user database.User) {

	feedFollows, err := s.DB.GetFeedFollowsForUser(r.Context(), user.ID)
	if err != nil {
		log.Printf("ERR during retriveal of all feedFollows: %v\n", err)
		respondWithError(w, http.StatusInternalServerError,"can't retrieve feedfollows")
		return
	}
	// IF NOT entries in feed-follows for this user?
	respondWithJSON(w, http.StatusOK, feedFollows)
	return
}

