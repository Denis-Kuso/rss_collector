package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"context"
	"io"
	"github.com/Denis-Kuso/rss_aggregator_p/internal/database"
	"github.com/google/uuid"
)

func (s *stateConfig) CreateUser(w http.ResponseWriter, r *http.Request){

	ctx := context.Background()
	type userRequest struct{
	Name string `json:"name"`
	}
	// parse request
	data, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest,"")// TODO better response
		return
	}
	userReq := userRequest{}
	err = json.Unmarshal(data, &userReq)
	if err != nil{
		respondWithError(w,http.StatusInternalServerError,"sorry")
		return
	}

	userId := uuid.New()
	user, err := s.DB.CreateUser(ctx, database.CreateUserParams{
		ID: userId,
		CreatedAt:time.Now(),
		UpdatedAt: time.Now(),// or save time.Now()?
		Name: userReq.Name,
	})
	if err != nil {
		log.Printf("Handle err:%v", err)
		respondWithError(w, http.StatusInternalServerError,"Sorry pal, cant make you")
		return
	}
	log.Printf("Created user: %v\n", user)
	respondWithJSON(w, http.StatusOK, user)
	return
}
