package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jbeyer16/boot-dev-rss-feed-agg/src/internal/database"
)

func (cfg *apiConfig) CreateUser(w http.ResponseWriter, r *http.Request) {
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

	user, err := cfg.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
		Name:      requestBody.Name,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Unable to create user (%v)", err))
		return
	}

	type response struct {
		Id        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Name      string    `json:"name"`
	}

	respondWithJSON(w, http.StatusCreated, response{
		Id:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
	})
}
