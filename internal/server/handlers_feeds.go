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

	type userFeedFollowResponse struct {
		Id        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		FeedId    uuid.UUID `json:"feed_id"`
		UserId    uuid.UUID `json:"user_id"`
	}
	type feedResponse struct {
		Id        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Name      string    `json:"name"`
		Url       string    `json:"url"`
		UserId    uuid.UUID `json:"user_id"`
	}
	response := struct {
		Feed        feedResponse           `json:"feed"`
		Feed_follow userFeedFollowResponse `json:"feed_follow"`
	}{
		Feed: feedResponse{
			Id:        feed.ID,
			CreatedAt: feed.CreatedAt,
			UpdatedAt: feed.UpdatedAt,
			Name:      feed.Name,
			Url:       feed.Url,
			UserId:    feed.UserID,
		},
		Feed_follow: userFeedFollowResponse{
			Id:        feedFollow.ID,
			CreatedAt: feedFollow.CreatedAt,
			UpdatedAt: feedFollow.UpdatedAt,
			FeedId:    feedFollow.FeedID,
			UserId:    feedFollow.UserID,
		},
	}
	respondWithJSON(w, http.StatusCreated, response)
}

func (s *Server) getFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := s.db.GetFeeds(r.Context())

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
