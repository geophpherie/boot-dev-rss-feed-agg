package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/geophpherie/boot-dev-rss-feed-agg/internal/database"
	"github.com/google/uuid"
)

func (s *Server) createFeed(w http.ResponseWriter, r *http.Request, user database.User) {
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

	feed, err := s.db.CreateFeed(r.Context(), database.CreateFeedParams{
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

	feedFollow, err := s.db.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
		FeedID:    feed.ID,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create feed follow")
		return
	}

	response := struct {
		Feed       Feed       `json:"feed"`
		FeedFollow FeedFollow `json:"feed_follow"`
	}{
		Feed:       ConvertFeedModelToResource(feed),
		FeedFollow: ConvertFeedFollowModelToResource(feedFollow),
	}
	respondWithJSON(w, http.StatusCreated, response)
}

func (s *Server) getFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := s.db.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to fetch feeds")
		return
	}

	response := make(Feeds, 0, len(feeds))
	for _, feed := range feeds {
		response = append(response, ConvertFeedModelToResource(feed))
	}

	respondWithJSON(w, http.StatusCreated, response)
}
