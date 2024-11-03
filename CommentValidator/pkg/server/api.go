package server

import (
	"commentvalidator/pkg/validate"
	"errors"
	"net/http"
)

var missingText = errors.New("missing comment text")

func (s *Server) Validate(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("text") {
		http.Error(w, "missing text parameter", http.StatusBadRequest)
		return
	}

	text := r.URL.Query().Get("text")
	if validate.IsValid(text) {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
