package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/Denis-Kuso/rss_aggregator_p/internal/database"
	"github.com/google/uuid"
)

func (s *StateConfig) GetAllFollowedFeeds(w http.ResponseWriter, r *http.Request, user database.User) {

	var errMsg string
	feedFollows, err := s.DB.GetFeedFollowsForUser(r.Context(), user.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errMsg = fmt.Sprintf("no followed feeds found")
			respondWithJSON(w, http.StatusOK, errMsg)
			return
		}
		errMsg = fmt.Sprintf("err retrieving feedFollows: %v", err)
		respondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}
	feedIDs := make([]uuid.UUID, len(feedFollows))
	for i, f := range feedFollows {
		feedIDs[i] = f.FeedID
	}
	feeds, err := s.DB.GetBasicInfoFeed(r.Context(), feedIDs)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			errMsg := fmt.Sprintf("cannot retrieve info about feeds: %s", err)
			respondWithError(w, http.StatusInternalServerError, errMsg)
			return
		}
	}
	publicFeeds := dbFeedToPublicFeeds(feeds)
	respondWithJSON(w, http.StatusOK, publicFeeds)
	return
}
