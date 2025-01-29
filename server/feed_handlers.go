package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Denis-Kuso/rss_collector/server/internal/database"
	"github.com/Denis-Kuso/rss_collector/server/internal/storage"
	"github.com/Denis-Kuso/rss_collector/server/internal/validate"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func (a *app) CreateFeed(w http.ResponseWriter, r *http.Request) {

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
			errMsg = fmt.Sprintf("cannot parse json, err at position: %d", jsonErr.Offset)
			respondWithError(w, http.StatusBadRequest, errMsg)
			return
		}

		a.serverErrorResponse(w, r, err)
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

	userID := r.Context().Value("userID").(uuid.UUID) // TODO generate-type-safe key as it stands this could panic
	err = a.feeds.Create(r.Context(), userID, userReq.Name, userReq.URL)
	if err != nil {
		if errors.Is(err, storage.ErrDuplicate) {
			errMsg = fmt.Sprintf("already following or you can follow %s", userReq.URL)
			respondWithError(w, http.StatusConflict, errMsg)
			return
		}
		err = fmt.Errorf("failed creating feed: %q: %q: %v", userReq.Name, userReq.URL, err)
		a.serverErrorResponse(w, r, err)
		return
	}
	respondWithJSON(w, http.StatusCreated, nil)
	return
}

func (a *app) FollowFeed(w http.ResponseWriter, r *http.Request) {
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
			errMsg = fmt.Sprintf("cannot parse json, err at position: %d", jsonErr.Offset)
			respondWithError(w, http.StatusBadRequest, errMsg)
			return
		}

		a.serverErrorResponse(w, r, err)
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
		err = fmt.Errorf("cannot follow feedID: %q: %v", userReq.FeedID, err)
		a.serverErrorResponse(w, r, err)
		return
	}
	userID := r.Context().Value("userID").(uuid.UUID) // TODO generate-type-safe key as it stands this could panic
	_, err = a.db.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    userID,
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

		err = fmt.Errorf("cannot follow feedID: %q: %v", userReq.FeedID, err)
		a.serverErrorResponse(w, r, err)
		return
	}
	pubFeed := dbFeedToPublicFeed(feedsInfo[0]) // use first and only element
	respondWithJSON(w, http.StatusOK, pubFeed)
}

func (a *app) GetAllFollowedFeeds(w http.ResponseWriter, r *http.Request) {

	var errMsg string
	userID := r.Context().Value("userID").(uuid.UUID) // TODO generate-type-safe key as it stands this could panic
	feedFollows, err := a.db.GetFeedFollowsForUser(r.Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errMsg = fmt.Sprintf("no followed feeds found")
			respondWithJSON(w, http.StatusOK, errMsg)
			return
		}
		err = fmt.Errorf("cannot retrieve feedfollows for user: %d: %v", userID, err)
		a.serverErrorResponse(w, r, err)
		return
	}
	feedIDs := make([]uuid.UUID, len(feedFollows))
	for i, f := range feedFollows {
		feedIDs[i] = f.FeedID
	}
	feeds, err := a.db.GetBasicInfoFeed(r.Context(), feedIDs)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		err = fmt.Errorf("cannot retrieve feed info: %d: %v", userID, err)
		a.serverErrorResponse(w, r, err)
		return
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
		err = fmt.Errorf("cannot retrieve feeds: %v", err)
		a.serverErrorResponse(w, r, err)
		return
	}
	publicFeeds := dbFeedToPublicFeeds(feeds)
	respondWithJSON(w, http.StatusOK, publicFeeds)
	return
}

func (a *app) UnfollowFeed(w http.ResponseWriter, r *http.Request) {
	var errMsg string
	type response struct {
		Name string `json:"unfollowedFeed"`
	}
	const queryKey = "feedFollowID"
	queries := r.URL.Query()
	providedFeedID := queries.Get(queryKey)

	feedID, err := uuid.Parse(providedFeedID)
	if err != nil {
		errMsg = fmt.Sprintf("Cannot parse feed id: %s", providedFeedID)
		respondWithError(w, http.StatusBadRequest, errMsg)
		return
	}

	userID := r.Context().Value("userID").(uuid.UUID) // TODO generate-type-safe key as it stands this could panic

	err = a.db.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		FeedID: feedID,
		UserID: userID,
	})
	if err != nil {
		err = fmt.Errorf("cannot delete following: %v: %v", feedID, err)
		a.serverErrorResponse(w, r, err)
		return
	}
	respondWithJSON(w, http.StatusOK, response{providedFeedID})
	return
}
