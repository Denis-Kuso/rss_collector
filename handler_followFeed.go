package main

import (
	"encoding/json"
	"fmt"
	"github.com/Denis-Kuso/rss_aggregator_p/internal/database"
	"github.com/google/uuid"
	"io"
	"net/http"
	"time"
)

func (s *stateConfig) FollowFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type userRequest struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	var errMsg string
	data, err := io.ReadAll(r.Body)
	if err != nil {
		errMsg = "err during reading response body"
		respondWithError(w, http.StatusBadRequest, errMsg)
		return
	}
	userReq := userRequest{}
	err = json.Unmarshal(data, &userReq)
	if err != nil {
		if jsonErr, ok := err.(*json.SyntaxError); ok {
			errMsg = fmt.Sprintf("cannot parse json, err occured at byte:%d", jsonErr.Offset)
			respondWithError(w, http.StatusBadRequest, errMsg)
			return
		}

		errMsg = "cannot parse json"
		respondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}
	feedFollow, err := s.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    userReq.FeedID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		errMsg = fmt.Sprintf("cannot follow feed with id: %s, err: %v", userReq.FeedID, err)
		respondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}
	respondWithJSON(w, http.StatusOK, feedFollow)
}
