package main

import (
	"log"
	"net/http"
	"github.com/Denis-Kuso/rss_aggregator_p/internal/database"
)

func (s *stateConfig) GetUserData(w http.ResponseWriter, r *http.Request, user database.User){

	log.Printf("Succesful retriveal of user data %v\n", user)
	respondWithJSON(w, http.StatusOK, user)
	return
}
