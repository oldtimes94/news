package server

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Server struct {
	Router *mux.Router
}

func New() *Server {
	return &Server{
		Router: mux.NewRouter(),
	}
}

func (s *Server) Register() {
	s.Router.HandleFunc("/news", s.news).Methods(http.MethodGet)
	s.Router.HandleFunc("/news/{id}", s.newsDetail).Methods(http.MethodGet)
	s.Router.HandleFunc("/comment", s.comment).Methods(http.MethodPost)
}
