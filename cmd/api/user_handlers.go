package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Denis-Kuso/rss_collector/internal/storage"
	"github.com/Denis-Kuso/rss_collector/internal/validate"
)

// TODO provide differently
const (
	queryKey          = "limit"
	defaultQueryLimit = 5
	maxPosts          = 100
)

func (a *app) CreateUser(w http.ResponseWriter, r *http.Request) {

	type userRequest struct {
		Name string `json:"name"`
	}
	userReq := userRequest{}
	err := readJSON(r, &userReq)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if userReq.Name == "" {
		respondWithError(w, http.StatusBadRequest, "body must not be empty")
		return
	}
	if ok := validate.ValidateUsername(userReq.Name); !ok {
		errMsg := fmt.Sprintf("invalid username: %s", userReq.Name)
		respondWithError(w, http.StatusBadRequest, errMsg)
		return
	}

	user, err := a.users.Create(r.Context(), userReq.Name)
	if err != nil {
		// TODO check for duplicate
		err = fmt.Errorf("cannot create user: %v: %v", userReq.Name, err)
		a.serverErrorResponse(w, r, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, user)
	return
}

func (a *app) GetPostsFromUser(w http.ResponseWriter, r *http.Request) {
	limit := defaultQueryLimit
	var errMsg string
	desiredLimit := r.URL.Query().Get(queryKey)
	// is limit parameter provided and smaller than max?
	if desiredLimit != "" {
		dLimit, err := strconv.Atoi(desiredLimit)
		if err != nil {
			errMsg = fmt.Sprintf("Provided limit value: %s not supported", desiredLimit)
			respondWithError(w, http.StatusBadRequest, errMsg)
			return
		}
		if (0 < dLimit) && (dLimit < maxPosts) {
			limit = dLimit
		}
	}
	userID, ok := GetUserIDFromContext(r)
	if !ok {
		slog.Warn("BUG - missing/empty userID", "userID", userID) // TODO here is where reqID might be useful
		respondWithError(w, http.StatusUnauthorized, "you must be authenticated to access this resource")
		return
	}
	posts, err := a.posts.Get(r.Context(), userID, limit)
	if err != nil {
		if !errors.Is(err, storage.ErrNotFound) {
			a.serverErrorResponse(w, r, err)
			return
		}
		// ok to bubble down (empty slice of posts is ok)
	}
	respondWithJSON(w, http.StatusOK, posts)
}

func (a *app) GetUserData(w http.ResponseWriter, r *http.Request) {

	userID, ok := GetUserIDFromContext(r)
	if !ok {
		slog.Warn("BUG - missing/empty userID", "userID", userID) // TODO here is where reqID might be useful
		respondWithError(w, http.StatusUnauthorized, "you must be authenticated to access this resource")
		return
	}
	u, err := a.users.Get(r.Context(), userID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) { // this should not really happen
			a.logError(r, err)
			respondWithError(w, http.StatusNotFound, "not found")
			return
		}
		a.serverErrorResponse(w, r, err)
		return
	}
	respondWithJSON(w, http.StatusOK, u)
	return
}
