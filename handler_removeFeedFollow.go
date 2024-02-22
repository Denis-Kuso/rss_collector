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
		respondWithError(w, http.StatusInternalServerError,"Cannot proces feed id")
		return
	}
	// Does feed_follow even exist?
	_, err = s.DB.GetFeedFollow(r.Context(), feedFollowID)
	if err != nil {
		log.Printf("%v. FeedFollow does not exists, cannot delete feedFollow :%v.\n",err, providedFeedFollowID)
		respondWithError(w, http.StatusBadRequest,"No feed follow for this user and feed\n")
		return
	}
	
	// err would be nil on non-existing entry
	err = s.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID: feedFollowID,
		UserID: user.ID,
	})
	if err != nil {
		log.Printf("ERR during removal of Feed-followd from db: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, "cannot remove link to feed")
		return
	}
	respondWithJSON(w, http.StatusOK,"Unfollowed feed")
	return
}
