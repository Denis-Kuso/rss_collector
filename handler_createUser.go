package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"fmt"
	"io"
	"github.com/Denis-Kuso/rss_aggregator_p/internal/database"
	"github.com/google/uuid"
)

func (s *stateConfig) CreateUser(w http.ResponseWriter, r *http.Request){

	type userRequest struct{
	Name string `json:"name"`
	}
	var errMsg string
	data, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest,"")// TODO better ERR handling 
		return
	}
	userReq := userRequest{}
	err = json.Unmarshal(data, &userReq)
	if err != nil{
		if jsonErr, ok := err.(*json.SyntaxError); ok {
			errMsg = fmt.Sprintf("cannot parse json, err occured at byte:%d", jsonErr.Offset)
			respondWithError(w, http.StatusBadRequest, errMsg)
			return
		}
		
		errMsg = "cannot parse json"
		respondWithError(w,http.StatusInternalServerError,errMsg)
		return
	}

	user, err := s.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt:time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: userReq.Name,
	})
	if err != nil {
		errMsg = fmt.Sprintf("cannot create user: %s", userReq.Name)
		log.Printf("%s, err: %v", errMsg, err)
		respondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}
	log.Printf("Created user: %v\n", user)
	respondWithJSON(w, http.StatusOK, user)
	return
}
