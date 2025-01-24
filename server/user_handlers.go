package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/Denis-Kuso/rss_collector/server/internal/database"
	"github.com/Denis-Kuso/rss_collector/server/internal/validate"
	"github.com/google/uuid"
)

const (
	QUERY_LIMIT         = "limit"
	DEFAULT_QUERY_LIMIT = 5
	MAX_PROVIDED_POSTS  = 100
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

	user, err := a.db.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      userReq.Name,
	})
	if err != nil {
		err = fmt.Errorf("cannot create user: %v: %v", userReq.Name, err)
		a.serverErrorResponse(w, r, err)
		return
	}

	publicUser := dbUserToPublicUser(user, make([]database.Feed, 0)) // no feeds for a new user
	respondWithJSON(w, http.StatusCreated, publicUser)
	return
}

func (a *app) GetPostsFromUser(w http.ResponseWriter, r *http.Request, user database.User) {
	limit := DEFAULT_QUERY_LIMIT
	var errMsg string
	desired_limit := r.URL.Query().Get(QUERY_LIMIT)
	// is limit parameter provided and smaller than max?
	if desired_limit != "" {
		desired_limit_I, err := strconv.Atoi(desired_limit)
		if err != nil {
			errMsg = fmt.Sprintf("Provided limit value: %s not supported", desired_limit)
			respondWithError(w, http.StatusBadRequest, errMsg)
			return
		}
		if (0 < desired_limit_I) && (desired_limit_I < MAX_PROVIDED_POSTS) {
			limit = desired_limit_I
		}
	}
	posts, err := a.db.GetPostsFromUser(r.Context(), database.GetPostsFromUserParams{
		UserID: user.ID,
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

func (a *app) GetUserData(w http.ResponseWriter, r *http.Request, user database.User) {

	feedFollows, err := a.db.GetFeedFollowsForUser(r.Context(), user.ID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			err = fmt.Errorf("cannot retrieve user info: %q: %v", user.ID, err)
			a.serverErrorResponse(w, r, err)
			return
		}
	}
	SIZE := len(feedFollows)
	feedIDs := make([]uuid.UUID, SIZE)
	for i, f := range feedFollows {
		feedIDs[i] = f.FeedID
	}
	feeds := make([]database.Feed, SIZE)
	feeds, err = a.db.GetBasicInfoFeed(r.Context(), feedIDs)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			err = fmt.Errorf("cannot retrieve feed info: %v", err)
			a.serverErrorResponse(w, r, err)
			return
		}
	}
	publicUser := dbUserToPublicUser(user, feeds)
	respondWithJSON(w, http.StatusOK, publicUser)
	return
}
