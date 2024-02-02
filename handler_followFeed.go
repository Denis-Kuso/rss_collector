package main

import (
	"encoding/json"
	"net/http"
	"io"
	"time"
	"github.com/Denis-Kuso/rss_aggregator_p/internal/database"
	"github.com/google/uuid"
)
func (s *stateConfig) FollowFeed(w http.ResponseWriter, r *http.Request, user database.User){
	type userRequest struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Check your request man")
		return}
	userReq := userRequest{}
	err = json.Unmarshal(data, &userReq)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse your request")
		return
	}
	// a feedfollow can be created for an existing feed-not merely when a feed is created
	// Should I also create a FEED if it does not exist?
	feedFollow, err := s.DB.CreateFeedFollow(r.Context(),database.CreateFeedFollowParams{
		ID: uuid.New(),
		UserID: user.ID,
		FeedID: userReq.FeedID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Can't create feed-follow")
		return
	}
	respondWithJSON(w, http.StatusOK, feedFollow)
}
