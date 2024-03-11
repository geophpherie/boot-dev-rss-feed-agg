package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/geophpherie/boot-dev-rss-feed-agg/src/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) createFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type body struct {
		Name string
		Url  string
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

	feed, err := cfg.db.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
		Name:      requestBody.Name,
		Url:       requestBody.Url,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create feed")
		return
	}

	type response struct {
		Id        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Name      string    `json:"name"`
		Url       string    `json:"url"`
		UserId    uuid.UUID `json:"user_id"`
	}
	respondWithJSON(w, http.StatusCreated, response{
		Id:        feed.ID,
		CreatedAt: feed.CreatedAt,
		UpdatedAt: feed.UpdatedAt,
		Name:      feed.Name,
		Url:       feed.Url,
		UserId:    feed.UserID,
	})
}

func (cfg *apiConfig) getFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := cfg.db.GetFeeds(r.Context())

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to fetch feeds")
		return
	}

	type feedResponse struct {
		Id        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Name      string    `json:"name"`
		Url       string    `json:"url"`
		UserId    uuid.UUID `json:"user_id"`
	}

	response := []feedResponse{}
	for _, feed := range feeds {
		response = append(response, feedResponse{
			Id:        feed.ID,
			CreatedAt: feed.CreatedAt,
			UpdatedAt: feed.UpdatedAt,
			Name:      feed.Name,
			Url:       feed.Url,
			UserId:    feed.UserID,
		})

	}
	respondWithJSON(w, http.StatusCreated, response)
}
