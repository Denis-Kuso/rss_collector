package main

import (
	"log"
	"net/http"
	"fmt"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/Denis-Kuso/rss_aggregator_p/internal/database"
)

func (s *stateConfig) GetUserData(w http.ResponseWriter, r *http.Request, user database.User){

	var errMsg string
	feedFollows, err := s.DB.GetFeedFollowsForUser(r.Context(), user.ID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
		errMsg = "err retrieving feed follows"
		respondWithError(w, http.StatusInternalServerError, "can't retrieve feedfollows")
		return
		}
	}
	SIZE := len(feedFollows)
	feedIDs := make([]uuid.UUID, SIZE)
	for i, f := range feedFollows{
		feedIDs[i] =  f.FeedID
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
	publicUser := dbUserToPublicUser(user,feeds)
	respondWithJSON(w, http.StatusOK, publicUser)
	return
}
