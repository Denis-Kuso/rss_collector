package main

import (
	"log"
	"net/http"
	"github.com/Denis-Kuso/rss_aggregator_p/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (s *stateConfig) UnfollowFeed(w http.ResponseWriter, r *http.Request, user database.User){
	providedFeedFollowID := chi.URLParam(r, QUERY_FEED_FOLLOW)

	feedFollowID, err := uuid.Parse(providedFeedFollowID)
	if err != nil{
		log.Printf("ERR during Feed-Follow id converstion to UUID: %v\n", err)
		respondWithJSON(w, http.StatusInternalServerError,"Cannot proces feed id")
		return
	}
	_, err = s.DB.RemoveFeedFollow(r.Context(), feedFollowID)
	if err != nil {
		log.Printf("ERR during removal of Feed-followd from db: %v\n", err)
		respondWithJSON(w, http.StatusInternalServerError, "cannot remove link to feed")
		return
	}
	respondWithJSON(w, http.StatusOK,"")
	return
}
