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

	type response struct {
		Id        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		FeedId    uuid.UUID `json:"feed_id"`
		UserId    uuid.UUID `json:"user_id"`
	}
	respondWithJSON(w, http.StatusCreated, response{
		Id:        feedFollow.ID,
		CreatedAt: feedFollow.CreatedAt,
		UpdatedAt: feedFollow.UpdatedAt,
		FeedId:    feedFollow.FeedID,
		UserId:    feedFollow.UserID,
	})
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

	type userFeedFollow struct {
		Id        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		FeedId    uuid.UUID `json:"feed_id"`
		UserId    uuid.UUID `json:"user_id"`
	}

	response := make([]userFeedFollow, 0, len(feedFollows))
	for _, feedFollow := range feedFollows {
		response = append(response, userFeedFollow{
			Id:        feedFollow.ID,
			CreatedAt: feedFollow.CreatedAt,
			UpdatedAt: feedFollow.UpdatedAt,
			FeedId:    feedFollow.FeedID,
			UserId:    feedFollow.UserID,
		})
	}
	respondWithJSON(w, http.StatusCreated, response)
}
