package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Denis-Kuso/rss_aggregator_p/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (s *stateConfig) UnfollowFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	var errMsg string
	providedFeedFollowID := chi.URLParam(r, QUERY_FEED_FOLLOW)

	feedFollowID, err := uuid.Parse(providedFeedFollowID)
	if err != nil {
		errMsg = fmt.Sprintf("Cannot parse feed id: %s", providedFeedFollowID)
		log.Printf("%s; err: %v\n", errMsg, err)
		respondWithError(w, http.StatusBadRequest, errMsg)
		return
	}

	err = s.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		FeedID: feedFollowID,
		UserID: user.ID,
	})
	if err != nil {
		errMsg = fmt.Sprintf("cannot unfollow feed: %s", providedFeedFollowID)
		log.Printf("%s; err: %v\n", errMsg, err)
		respondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}
	respondWithJSON(w, http.StatusOK, "Unfollowed feed")
	return
}
