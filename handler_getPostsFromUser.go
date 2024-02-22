package main

import (
	"log"
	"net/http"
	"github.com/Denis-Kuso/rss_aggregator_p/internal/database"
	"strconv"
)

const (
	QUERY_LIMIT = "limit"
	DEFAULT_QUERY_LIMIT = 5
	MAX_PROVIDED_POSTS = 100
)


func (s *stateConfig) GetPostsFromUser(w http.ResponseWriter, r *http.Request, user database.User) {
	limit := DEFAULT_QUERY_LIMIT
	desired_limit := r.URL.Query().Get(QUERY_LIMIT)
	// is limit parameter provided and smaller than max?
	if desired_limit != "" {
		desired_limit_I, err := strconv.Atoi(desired_limit)
		if err != nil {
			log.Printf("ERR during limit query conversion: %v\n", err)
			respondWithError(w, http.StatusBadRequest,"Provided value for limit not supported\n")
			return
		}
		if  (0 < desired_limit_I) && (desired_limit_I < MAX_PROVIDED_POSTS) {
			limit = desired_limit_I
		}
	}
	// in absence return default num of posts
	posts, err := s.DB.GetPostsFromUser(r.Context(), database.GetPostsFromUserParams{
		UserID: user.ID,
		Limit: int32(limit),
		})
	if err != nil {
		log.Printf("ERR during post retrieval: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, "Could not retrieve posts\n")
		return
	}
	respondWithJSON(w, http.StatusOK, posts)
}
