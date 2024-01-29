package main

import (
	"errors"
	"log"
	"io"
	"encoding/json"
	"net/http"
	"github.com/Denis-Kuso/rss_aggregator_p/internal/auth"
	"github.com/Denis-Kuso/rss_aggregator_p/internal/database"
	"time"
	"github.com/google/uuid"
)

func (s *stateConfig) CreateFeed(w http.ResponseWriter, r *http.Request){

	type Request struct {
		Name string `json:"name"`
		URL string `json:"url"` // PERHAPS URL TYPE?
	}

	// Check API key
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
	data, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest,"")// TODO better response
		return
	}
	userReq := Request{}
	err = json.Unmarshal(data, &userReq)
	if err != nil{
		respondWithError(w,http.StatusInternalServerError, "sorry")
		return
	}
	// SOME INPUT HANDLING ON NAME AND URL?
	user, err := s.DB.GetUserByAPI(r.Context(), apiKey)
	if err != nil{
		respondWithError(w, http.StatusUnauthorized, "NO USER INFO")
		return
	}

	feed, err := s.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID: user.ID,
		Name: userReq.Name,
		Url: userReq.URL,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError,"Internal err")
		log.Printf("ERR: %v\n",err)
		return
	}
	log.Printf("Succesful creation of feed %v\n", feed)
	respondWithJSON(w, http.StatusOK, feed)
	return
}
