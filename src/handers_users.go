package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/geophpherie/boot-dev-rss-feed-agg/src/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	type body struct {
		Name string
	}

	decoder := json.NewDecoder(r.Body)
	requestBody := body{}
	err := decoder.Decode(&requestBody)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to parse request body")
		return
	}

	id, err := uuid.NewV7()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to generate id")
		return
	}
	now := time.Now()

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
		Name:      requestBody.Name,
	})
	if err != nil {
		log.Printf("Unable to create user (%v)", err)
		respondWithError(w, http.StatusBadRequest, "Couldn't create user")
		return
	}

	type response struct {
		Id        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Name      string    `json:"name"`
		ApiKey    string    `json:"api_key"`
	}

	respondWithJSON(w, http.StatusCreated, response{
		Id:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		ApiKey:    user.ApiKey,
	})
}

func (cfg *apiConfig) getUser(w http.ResponseWriter, r *http.Request, user database.User) {
	type response struct {
		Id        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Name      string    `json:"name"`
		ApiKey    string    `json:"api_key"`
	}

	respondWithJSON(w, http.StatusOK, response{
		Id:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		ApiKey:    user.ApiKey,
	})
}
