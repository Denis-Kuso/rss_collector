package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	var errMsg string
	data, err := io.ReadAll(r.Body)
	if err != nil {
		errMsg = fmt.Sprintf("could not read request: %v", err)
		respondWithError(w, http.StatusBadRequest, errMsg)
		return
	}
	userReq := userRequest{}
	err = json.Unmarshal(data, &userReq)
	// TODO create custom JSON messages
	if err != nil {
		if jsonErr, ok := err.(*json.SyntaxError); ok {
			errMsg = fmt.Sprintf("cannot parse json, err occured at position: %d", jsonErr.Offset)
			respondWithError(w, http.StatusBadRequest, errMsg)
			return
		}

		errMsg = "cannot parse json"
		respondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}
	if ok := validate.ValidateUsername(userReq.Name); !ok {
		errMsg = fmt.Sprintf("invalid username: %s", userReq.Name)
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
	userID := r.Context().Value("userID").(uuid.UUID) // TODO generate-type-safe key as it stands this could panic
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

	userID := r.Context().Value("userID").(uuid.UUID) // TODO this can panic
	u, err := a.users.Get(r.Context(), userID)
	if err != nil {
		// TODO what could happen here
		if errors.Is(err, storage.ErrNotFound) { // TODO this should not really happen and it should be unauthorised
			respondWithError(w, http.StatusNotFound, "not found")
			return
		}
	}
	respondWithJSON(w, http.StatusOK, u)
	return
}
