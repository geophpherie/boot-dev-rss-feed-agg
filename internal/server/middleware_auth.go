package server

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/geophpherie/boot-dev-rss-feed-agg/internal/auth"
	"github.com/geophpherie/boot-dev-rss-feed-agg/internal/database"
)

type authenticatedHandler func(http.ResponseWriter, *http.Request, database.User)

func (s *Server) middlewareAuth(next authenticatedHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.ParseApiKey(r.Header)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "unable to parse api key")
			return
		}

		user, err := s.db.GetUser(r.Context(), apiKey)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				respondWithError(w, http.StatusUnauthorized, "api key not valid")
				return
			}
		}

		next(w, r, user)
	})
}
