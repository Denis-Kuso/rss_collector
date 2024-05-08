package main

import (
	"encoding/json"
	"fmt"
	"github.com/Denis-Kuso/rss_collector/server/internal/database"
	"github.com/Denis-Kuso/rss_collector/server/internal/validate"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"time"
)

func (s *StateConfig) CreateUser(w http.ResponseWriter, r *http.Request) {

	type userRequest struct {
		Name string `json:"name"`
	}
	var errMsg string
	data, err := io.ReadAll(r.Body)
	if err != nil {
		errMsg = fmt.Sprintf("could not read request: %v", err)
		respondWithError(w, http.StatusBadRequest, errMsg)
		return
	}
	userReq := userRequest{}
	err = json.Unmarshal(data, &userReq)
	if err != nil {
		if jsonErr, ok := err.(*json.SyntaxError); ok {
			errMsg = fmt.Sprintf("cannot parse json, err occured at byte:%d", jsonErr.Offset)
			respondWithError(w, http.StatusBadRequest, errMsg)
			return
		}

		errMsg = "cannot parse json"
		respondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}
	if ok := validate.ValidateUsername(userReq.Name); !ok {
		errMsg = fmt.Sprintf("invalid username: %s", userReq.Name)
		respondWithError(w, http.StatusBadRequest, errMsg)
		return
	}

	user, err := s.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      userReq.Name,
	})
	if err != nil {
		errMsg = fmt.Sprintf("cannot create user: %s", userReq.Name)
		log.Printf("%s, err: %v", errMsg, err)
		respondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}
	log.Printf("Created user: %v\n", user)

	publicUser := dbUserToPublicUser(user, make([]database.Feed, 0)) // no feeds for a new user
	respondWithJSON(w, http.StatusOK, publicUser)
	return
}
