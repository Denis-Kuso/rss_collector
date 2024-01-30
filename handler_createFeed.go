package main

import (
	"log"
	"io"
	"encoding/json"
	"net/http"
	"github.com/Denis-Kuso/rss_aggregator_p/internal/database"
	"time"
	"github.com/google/uuid"
)

func (s *stateConfig) CreateFeed(w http.ResponseWriter, r *http.Request, user database.User){

	type Request struct {
		Name string `json:"name"`
		URL string `json:"url"` // PERHAPS URL TYPE?
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
