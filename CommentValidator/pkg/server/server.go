package server

import (
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
)

type Server struct {
	Logger *zap.Logger
	Router *mux.Router
}

func New() (*Server, error) {

	s := new(Server)
	s.Router = mux.NewRouter()

	return s, nil
}

func (s *Server) Register() {
	s.Router.HandleFunc("/validate", s.Validate).Methods(http.MethodGet)
}
