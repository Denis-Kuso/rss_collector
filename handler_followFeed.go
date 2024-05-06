package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Denis-Kuso/rss_aggregator_p/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
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
	FeedID := []uuid.UUID{userReq.FeedID}
	// does desired feed even exist?
	feedsInfo, err := s.DB.GetBasicInfoFeed(r.Context(), FeedID)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			errMsg = fmt.Sprintf("cannot follow feed: %s. No such feed", userReq.FeedID)
			respondWithError(w, http.StatusNotFound, errMsg)
			return
		}
		errMsg = fmt.Sprintf("cannot follow feed:%s :%v", userReq.FeedID, err)
		respondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}
	_, err = s.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    userReq.FeedID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			// unique key violation https://www.postgresql.org/docs/current/errcodes-appendix.html
			if err.Code == "23505" {
				errMsg = fmt.Sprintf("already following feed: %s", userReq.FeedID)
				respondWithError(w, http.StatusBadRequest, errMsg)
				return
			}
		}

		errMsg = fmt.Sprintf("cannot follow feed: %s, err: %v", userReq.FeedID, err)
		respondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}
	pubFeed := dbFeedToPublicFeed(feedsInfo[0]) // use first and only element
	respondWithJSON(w, http.StatusOK, pubFeed)
}
