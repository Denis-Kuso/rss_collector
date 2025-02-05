package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/Denis-Kuso/rss_collector/server/internal/storage"
	"github.com/Denis-Kuso/rss_collector/server/internal/validate"
	"github.com/google/uuid"
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

	userID, ok := GetUserIDFromContext(r)
	if !ok {
		slog.Warn("BUG - missing/empty userID", "userID", userID) // TODO here is where reqID might be useful // could could logError or logWarning
		respondWithError(w, http.StatusUnauthorized, "you must be authenticated to access this resource")
		return
	}
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
	userID, ok := GetUserIDFromContext(r)
	if !ok {
		slog.Warn("BUG - missing/empty userID", "userID", userID) // TODO here is where reqID might be useful // could could logError or logWarning
		respondWithError(w, http.StatusUnauthorized, "you must be authenticated to access this resource")
		return
	}
	err = a.feeds.Follow(r.Context(), userID, userReq.FeedID)
	if err != nil {
		switch err {
		case storage.ErrNotFound:
			respondWithError(w, http.StatusNotFound, "resource not found")
			return
		case storage.ErrDuplicate:
			respondWithError(w, http.StatusConflict, "already following")
			return
		default:
			a.serverErrorResponse(w, r, err)
			return
		}
	}
	// TODO do I return anything (or empty body) or 204?
	respondWithJSON(w, http.StatusOK, nil)
}

func (a *app) GetAllFollowedFeeds(w http.ResponseWriter, r *http.Request) {

	userID, ok := GetUserIDFromContext(r)
	if !ok {
		slog.Warn("BUG - missing/empty userID", "userID", userID) // TODO here is where reqID might be useful // could could logError or logWarning
		respondWithError(w, http.StatusUnauthorized, "you must be authenticated to access this resource")
		return
	}
	feeds, err := a.feeds.Get(r.Context(), userID)
	if err != nil {
		if !errors.Is(err, storage.ErrNotFound) {
			a.serverErrorResponse(w, r, err)
			return
		}
		// fallthrough is ok, as we return empty []Feed (not following anything )
	}
	respondWithJSON(w, http.StatusOK, feeds)
	return
}

func (a *app) GetFeeds(w http.ResponseWriter, r *http.Request) {

	feeds, err := a.feeds.ShowAvailable(r.Context())
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			respondWithJSON(w, http.StatusOK, "no feeds")
			return
		}
		err = fmt.Errorf("cannot retrieve feeds: %v", err)
		a.serverErrorResponse(w, r, err)
		return
	}
	respondWithJSON(w, http.StatusOK, feeds)
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

	userID, ok := GetUserIDFromContext(r)
	if !ok {
		slog.Warn("BUG - missing/empty userID", "userID", userID) // TODO here is where reqID might be useful // could could logError or logWarning
		respondWithError(w, http.StatusUnauthorized, "you must be authenticated to access this resource")
		return
	}

	err = a.feeds.Delete(r.Context(), feedID, userID)
	if err != nil {
		err = fmt.Errorf("cannot delete following: %v: %v", feedID, err)
		a.serverErrorResponse(w, r, err)
		return
	}
	respondWithJSON(w, http.StatusOK, response{providedFeedID})
	return
}
