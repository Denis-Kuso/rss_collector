package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Denis-Kuso/rss_aggregator_p/internal/database"
)

const (
	QUERY_LIMIT         = "limit"
	DEFAULT_QUERY_LIMIT = 5
	MAX_PROVIDED_POSTS  = 100
)

func (s *stateConfig) GetPostsFromUser(w http.ResponseWriter, r *http.Request, user database.User) {
	limit := DEFAULT_QUERY_LIMIT
	var errMsg string
	desired_limit := r.URL.Query().Get(QUERY_LIMIT)
	// is limit parameter provided and smaller than max?
	if desired_limit != "" {
		desired_limit_I, err := strconv.Atoi(desired_limit)
		if err != nil {
			errMsg = fmt.Sprintf("Provided limit value: %s not supported", desired_limit)
			log.Printf("%s, %v", errMsg, err) // is this handling the error twice?
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
	respondWithJSON(w, http.StatusOK, posts)
}
