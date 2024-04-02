package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/geophpherie/boot-dev-rss-feed-agg/internal/database"
)

func (s *Server) getUserPosts(w http.ResponseWriter, r *http.Request, user database.User) {
	limitParam := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		limit = 10
	}

	params := database.GetPostsByUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	}

	fmt.Println("!!!!!!!!!!!! limit is ", int32(limit))
	userPosts, err := s.db.GetPostsByUser(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to fetch user posts")
		return
	}

	response := make([]Post, 0, len(userPosts))
	for _, userPost := range userPosts {
		response = append(response, ConvertPostModelToResource(userPost))
	}
	respondWithJSON(w, http.StatusOK, response)
}
