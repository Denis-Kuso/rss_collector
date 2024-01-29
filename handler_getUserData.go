package main

import (
	"errors"
	"log"
	"net/http"
	"github.com/Denis-Kuso/rss_aggregator_p/internal/auth"
)

func (s *stateConfig) GetUserData(w http.ResponseWriter, r *http.Request){

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil{
		if errors.Is(err, auth.ErrNoAuthHeaderIncluded){
			respondWithError(w, http.StatusBadRequest, "NO HEADER INCLUDED")
		}else{
			respondWithError(w, http.StatusUnauthorized,"ERR during processing apiKey")
		}
		log.Printf("ERR: %s\n", err)
		return
	}

	user, err := s.DB.GetUserByAPI(r.Context(), apiKey)
	if err != nil {
		log.Printf("Handle err:%v", err)
		// if no user, error should be not found, not internal server error
		respondWithError(w, http.StatusInternalServerError,"Sorry, no user data.")
		return
	}
	log.Printf("Succesful retriveal of user data %v\n", user)
	respondWithJSON(w, http.StatusOK, user)
	return
}
