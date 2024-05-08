package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Denis-Kuso/rss_collector/server/internal/database"
	"github.com/google/uuid"
)

const (
	QUERY_LIMIT         = "limit"
	DEFAULT_QUERY_LIMIT = 5
	MAX_PROVIDED_POSTS  = 100
)

func (s *StateConfig) GetPostsFromUser(w http.ResponseWriter, r *http.Request, user database.User) {
	limit := DEFAULT_QUERY_LIMIT
	var errMsg string
	desired_limit := r.URL.Query().Get(QUERY_LIMIT)
	// is limit parameter provided and smaller than max?
	if desired_limit != "" {
		desired_limit_I, err := strconv.Atoi(desired_limit)
		if err != nil {
			errMsg = fmt.Sprintf("Provided limit value: %s not supported", desired_limit)
			log.Printf("%s, %v", errMsg, err)
			respondWithError(w, http.StatusBadRequest, errMsg)
			return
		}
		if (0 < desired_limit_I) && (desired_limit_I < MAX_PROVIDED_POSTS) {
			limit = desired_limit_I
		}
	}
	posts, err := s.DB.GetPostsFromUser(r.Context(), database.GetPostsFromUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errMsg = "no posts found"
			respondWithJSON(w, http.StatusOK, errMsg)
			return
		}
		errMsg = "could not retrieve posts"
		log.Printf("%s; key: %s; err: %v\n", errMsg, user.ApiKey, err)
		respondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}
	SIZE := len(posts)
	const FIRST int = 0
	feedID := make([]uuid.UUID, 1) // need an array/slice for sql query
	feeds := make([]database.Feed, SIZE)
	for i, p := range posts {
		feedID[FIRST] = p.FeedID
		feed, err := s.DB.GetBasicInfoFeed(r.Context(), feedID)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				errMsg = fmt.Sprintf("cannot retrieve info. Feed id: %v, err:%v", feedID, err)
				respondWithError(w, http.StatusInternalServerError, errMsg)
				return
			}
			continue
		}
		feeds[i] = feed[FIRST]
	}
	publicPosts := dbPostsToPublicPosts(posts, feeds)
	respondWithJSON(w, http.StatusOK, publicPosts)
}
