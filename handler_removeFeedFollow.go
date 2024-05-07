package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Denis-Kuso/rss_aggregator_p/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (s *StateConfig) UnfollowFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	var errMsg string
	type response struct {
		Name string `json:"unfollowedFeed"`
	}
	const QUERY_FEED_FOLLOW = "feedFollowID"
	providedFeedID := chi.URLParam(r, QUERY_FEED_FOLLOW)

	feedID, err := uuid.Parse(providedFeedID)
	if err != nil {
		errMsg = fmt.Sprintf("Cannot parse feed id: %s", providedFeedID)
		log.Printf("%s; err: %v\n", errMsg, err)
		respondWithError(w, http.StatusBadRequest, errMsg)
		return
	}

	err = s.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		FeedID: feedID,
		UserID: user.ID,
	})
	if err != nil {
		errMsg = fmt.Sprintf("cannot unfollow feed: %s", providedFeedID)
		log.Printf("%s; err: %v\n", errMsg, err)
		respondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}
	respondWithJSON(w, http.StatusOK, response{providedFeedID})
	return
}
