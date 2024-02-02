package main 

import (
	//"encoding/json"
	"log"
	"net/http"
	//"github.com/Denis-Kuso/rss_aggregator_p/internal/database"
)

func (s *stateConfig) GetFeeds(w http.ResponseWriter, r * http.Request) {

	feeds, err := s.DB.GetFeeds(r.Context())
	if err != nil {
		// TODO create more error granularity (not found vs internal error)
		log.Printf(" [ ERR ] - GET ALL FEEDS:%v\n", err)
		respondWithError(w, http.StatusNotFound, "No feeds")
		return
	}
	respondWithJSON(w, http.StatusOK, feeds)
	return
}
