package main

import (
	"errors"
	"log"
	"net/http"
	"github.com/Denis-Kuso/rss_aggregator_p/internal/auth"
	"github.com/Denis-Kuso/rss_aggregator_p/internal/database"
)
type authenicatedHandler func(w http.ResponseWriter, r *http.Request, user database.User)
func (s *stateConfig) MiddlewareAuth(handler authenicatedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){

	// Check API key
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil{
		if errors.Is(err, auth.ErrNoAuthHeaderIncluded){
			respondWithError(w, http.StatusBadRequest, "NO HEADER INCLUDED")
			return
		}else{
			respondWithError(w, http.StatusUnauthorized,"ERR during processing apiKey")
			return
		}
	}
	// get user with authenticated api 
	user, err := s.DB.GetUserByAPI(r.Context(), apiKey)
	if err != nil {
		log.Printf("Handle err: %v", err)
		respondWithError(w, http.StatusNotFound, "Sorry, no user data.")
		return
	}
	handler(w, r, user)
	}
}

