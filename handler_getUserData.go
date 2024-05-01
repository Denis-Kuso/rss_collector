package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Denis-Kuso/rss_aggregator_p/internal/database"
	"github.com/google/uuid"
	"net/http"
)

func (s *stateConfig) GetUserData(w http.ResponseWriter, r *http.Request, user database.User) {

	var errMsg string
	feedFollows, err := s.DB.GetFeedFollowsForUser(r.Context(), user.ID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			errMsg = fmt.Sprintf("err retrieving feed follows: %v", err)
			respondWithError(w, http.StatusInternalServerError, errMsg)
			return
		}
	}
	SIZE := len(feedFollows)
	feedIDs := make([]uuid.UUID, SIZE)
	for i, f := range feedFollows {
		feedIDs[i] = f.FeedID
	}
	feeds := make([]database.Feed, SIZE)
	feeds, err = s.DB.GetBasicInfoFeed(r.Context(), feedIDs)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			errMsg = fmt.Sprintf("cannot retrieve feed info: %v", err)
			respondWithError(w, http.StatusInternalServerError, errMsg)
			return
		}
	}
	publicUser := dbUserToPublicUser(user, feeds)
	respondWithJSON(w, http.StatusOK, publicUser)
	return
}
