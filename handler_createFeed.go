package main

import (
	"log"
	"io"
	"encoding/json"
	"net/http"
	"github.com/Denis-Kuso/rss_aggregator_p/internal/database"
	"time"
	"fmt"
	"github.com/google/uuid"
)

func (s *stateConfig) CreateFeed(w http.ResponseWriter, r *http.Request, user database.User){

	type Request struct {
		Name string `json:"name"`
		URL string `json:"url"` 
	}
	var errMsg string
	data, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest,"")// TODO better response
		return
	}
	userReq := Request{}
	err = json.Unmarshal(data, &userReq)
	if err != nil{
		if jsonErr, ok := err.(*json.SyntaxError); ok {
			errMsg = fmt.Sprintf("cannot parse json, err occured at byte:%d", jsonErr.Offset)
			respondWithError(w, http.StatusBadRequest, errMsg)
			return
		}
		
		errMsg = "cannot parse json"
		respondWithError(w,http.StatusInternalServerError,errMsg)
		return
	}

	feed, err := s.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID: user.ID,
		Name: userReq.Name,
		Url: userReq.URL,
	})

	if err != nil {
		errMsg = fmt.Sprintf("cannot create a following to feed: %s; %s",userReq.Name, userReq.URL)
		log.Printf("failed during feed creation: %v, %s; %s\n", err, userReq.Name, userReq.URL)
		respondWithError(w, http.StatusInternalServerError,errMsg)
		return
	}
	_, err = s.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})

	if err != nil {
		log.Printf("failed during feed creation: %v, %s; %s\n", err, userReq.Name, userReq.URL)
		respondWithError(w, http.StatusInternalServerError,errMsg)
		return
	}
	publicFeed := dbFeedToPublicFeed(feed)
	respondWithJSON(w, http.StatusOK, publicFeed)
	return
}
