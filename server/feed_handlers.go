package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Denis-Kuso/rss_collector/server/internal/database"
	"github.com/Denis-Kuso/rss_collector/server/internal/validate"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func (a *app) CreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {

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

	feed, err := a.db.CreateFeed(r.Context(), database.CreateFeedParams{
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
	_, err = a.db.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
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

func (a *app) FollowFeed(w http.ResponseWriter, r *http.Request, user database.User) {
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
	feedsInfo, err := a.db.GetBasicInfoFeed(r.Context(), FeedID)
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
	_, err = a.db.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
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

func (a *app) GetAllFollowedFeeds(w http.ResponseWriter, r *http.Request, user database.User) {

	var errMsg string
	feedFollows, err := a.db.GetFeedFollowsForUser(r.Context(), user.ID)
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
	feeds, err := a.db.GetBasicInfoFeed(r.Context(), feedIDs)
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

func (a *app) GetFeeds(w http.ResponseWriter, r *http.Request) {

	var errMsg string
	feeds, err := a.db.GetFeeds(r.Context())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errMsg = "no feeds found"
			respondWithJSON(w, http.StatusOK, errMsg)
			return
		}
		errMsg = "could not retrieve feeds"
		log.Printf("%s: %v\n", errMsg, err)
		respondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}
	publicFeeds := dbFeedToPublicFeeds(feeds)
	respondWithJSON(w, http.StatusOK, publicFeeds)
	return
}

func (a *app) UnfollowFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	var errMsg string
	type response struct {
		Name string `json:"unfollowedFeed"`
	}
	const QUERY_FEED_FOLLOW = "feedFollowID"
	providedFeedID := chi.URLParam(r, QUERY_FEED_FOLLOW)

	feedID, err := uuid.Parse(providedFeedID)
	if err != nil {
		errMsg = fmt.Sprintf("Cannot parse feed id: %s", providedFeedID)
		log.Printf("%s; err: %v\n", errMsg, err)
		respondWithError(w, http.StatusBadRequest, errMsg)
		return
	}

	err = a.db.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		FeedID: feedID,
		UserID: user.ID,
	})
	if err != nil {
		errMsg = fmt.Sprintf("cannot unfollow feed: %s", providedFeedID)
		log.Printf("%s; err: %v\n", errMsg, err)
		respondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}
	respondWithJSON(w, http.StatusOK, response{providedFeedID})
	return
}
