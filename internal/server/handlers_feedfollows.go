package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/geophpherie/boot-dev-rss-feed-agg/internal/database"
	"github.com/google/uuid"
)

func (s *Server) createFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type body struct {
		Feed_id string
	}

	decoder := json.NewDecoder(r.Body)
	requestBody := body{}
	err := decoder.Decode(&requestBody)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to parse request body")
		return
	}

	feedId, err := uuid.Parse(requestBody.Feed_id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid feedId format")
		return
	}

	id, err := uuid.NewV7()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to generate id")
	}

	now := time.Now()

	feedFollow, err := s.db.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
		FeedID:    feedId,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create feed follow")
		return
	}

	respondWithJSON(w, http.StatusCreated, ConvertFeedFollowModelToResource(feedFollow))
}

func (s *Server) deleteFeedFollow(w http.ResponseWriter, r *http.Request) {
	feedFollowId, err := uuid.Parse(r.PathValue("feed_follow_id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid feedFollowId format")
		return
	}

	err = s.db.DeleteFeedFollow(r.Context(), feedFollowId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to delete feed follow")
		return
	}

	respondWithJSON(w, http.StatusNoContent, struct{}{})
}

func (s *Server) getFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := s.db.GetUserFeedFollows(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to fetch user feed")
		return
	}

	response := make(FeedFollows, 0, len(feedFollows))
	for _, feedFollow := range feedFollows {
		response = append(response, ConvertFeedFollowModelToResource(feedFollow))
	}
	respondWithJSON(w, http.StatusCreated, response)
}
