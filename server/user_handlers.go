package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Denis-Kuso/rss_collector/server/internal/database"
	"github.com/Denis-Kuso/rss_collector/server/internal/storage"
	"github.com/Denis-Kuso/rss_collector/server/internal/validate"
	"github.com/google/uuid"
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
	posts, err := a.db.GetPostsFromUser(r.Context(), database.GetPostsFromUserParams{
		UserID: userID,
		Limit:  int32(limit),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errMsg = "no posts found"
			respondWithJSON(w, http.StatusOK, errMsg)
			return
		}
		err = fmt.Errorf("cannot retrieve posts: %v", err)
		a.serverErrorResponse(w, r, err)
		return
	}
	SIZE := len(posts)
	const FIRST int = 0
	feedID := make([]uuid.UUID, 1) // need an array/slice for sql query
	feeds := make([]database.Feed, SIZE)
	for i, p := range posts {
		feedID[FIRST] = p.FeedID
		feed, err := a.db.GetBasicInfoFeed(r.Context(), feedID)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				err = fmt.Errorf("cannot retrieve feed info: %q: %v", p.FeedID, err)
				a.serverErrorResponse(w, r, err)
				return
			}
			continue
		}
		feeds[i] = feed[FIRST]
	}
	publicPosts := dbPostsToPublicPosts(posts, feeds)
	respondWithJSON(w, http.StatusOK, publicPosts)
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
