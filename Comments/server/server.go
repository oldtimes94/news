package server

import (
	"comments/storage"
	"comments/storage/sqlite"
	"fmt"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
)

var Logger *zap.Logger

type Server struct {
	DB     storage.Interface
	Router *mux.Router
}

func New(storageType string) (*Server, error) {

	s := new(Server)
	s.Router = mux.NewRouter()

	switch storageType {
	case "sqlite":
		s.DB = sqlite.New()
	default:
		return nil, fmt.Errorf("unknown storage type: %s", storageType)
	}
	return s, nil
}

func (s *Server) Register() {
	s.Router.HandleFunc("/news", s.Post).Methods(http.MethodPost)
	s.Router.HandleFunc("/news/{id:[0-9]+}", s.CommentByNewsID).Methods(http.MethodGet)
}
