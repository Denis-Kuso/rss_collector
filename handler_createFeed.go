package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Denis-Kuso/rss_aggregator_p/internal/database"
	"github.com/Denis-Kuso/rss_aggregator_p/internal/validate"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func (s *stateConfig) CreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {

	type Request struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	var errMsg string
	data, err := io.ReadAll(r.Body)
	if err != nil {
		errMsg = fmt.Sprintf("cannot read request body: %v\n", err)
		respondWithError(w, http.StatusBadRequest, errMsg)
		return
	}
	userReq := Request{}
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
	if ok := validate.IsUrl(userReq.URL); !ok {
		errMsg = fmt.Sprintf("invalid url format: %s", userReq.URL)
		respondWithError(w, http.StatusBadRequest, errMsg)
		return
	}
	if ok := validate.ValidateUsername(userReq.Name); !ok { // same checks as username for now
		errMsg = fmt.Sprintf("invalid feed name: %s", userReq.Name)
		respondWithError(w, http.StatusBadRequest, errMsg)
		return
	}

	feed, err := s.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		Name:      userReq.Name,
		Url:       userReq.URL,
	})
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			// unique key violation https://www.postgresql.org/docs/current/errcodes-appendix.html
			if err.Code == "23505" {
				errMsg = fmt.Sprintf("\"%s\" already exists, try following the feed instead", userReq.Name)
				respondWithError(w, http.StatusConflict, errMsg)
				return
			}
		}
		errMsg = fmt.Sprintf("cannot create a following to feed: %s; %s", userReq.Name, userReq.URL)
		log.Printf("failed during feed creation: %v, %s; %s\n", err, userReq.Name, userReq.URL)
		respondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}
	_, err = s.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})

	if err != nil {
		log.Printf("failed during feed creation: %v, %s; %s\n", err, userReq.Name, userReq.URL)
		respondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}
	publicFeed := dbFeedToPublicFeed(feed)
	respondWithJSON(w, http.StatusOK, publicFeed)
	return
}
