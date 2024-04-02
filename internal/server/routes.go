package server

import (
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {

	serveMux := http.NewServeMux()

	serveMux.HandleFunc("GET /v1/readiness", handleReadiness)
	serveMux.HandleFunc("GET /v1/error", handleError)

	serveMux.HandleFunc("POST /v1/users", s.createUser)
	serveMux.HandleFunc("GET /v1/users", s.middlewareAuth(s.getUser))

	serveMux.HandleFunc("POST /v1/feeds", s.middlewareAuth(s.createFeed))
	serveMux.HandleFunc("GET /v1/feeds", s.getFeeds)

	serveMux.HandleFunc("POST /v1/feed_follows", s.middlewareAuth(s.createFeedFollow))
	serveMux.HandleFunc("GET /v1/feed_follows", s.middlewareAuth(s.getFeedFollows))
	serveMux.HandleFunc("DELETE /v1/feed_follows/{feed_follow_id}", s.deleteFeedFollow)

	serveMux.HandleFunc("GET /v1/posts", s.middlewareAuth(s.getUserPosts))
	return middlewareCors(serveMux)
}
