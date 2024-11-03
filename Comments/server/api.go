package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func (s *Server) Post(w http.ResponseWriter, r *http.Request) {
	newsID := r.URL.Query().Get("news_id")
	parentCommentID := r.URL.Query().Get("parent_comment_id")
	text := r.URL.Query().Get("text")

	if newsID == "" || text == "" {
		http.Error(w, "invalid input", http.StatusBadRequest)
		log.Println("invalid input")
		return
	}

	newsIDInt, err := strconv.Atoi(newsID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	parentCommentIDInt := 0
	if parentCommentID != "" {
		parentCommentIDInt, err = strconv.Atoi(parentCommentID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	err = s.DB.Add(newsIDInt, parentCommentIDInt, text)
	if err != nil {
		http.Error(w, "failed to add comment", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Comment added successfully"))
}

func (s *Server) CommentByNewsID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	newsID := vars["id"]

	newsIDInt, err := strconv.Atoi(newsID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	comments, err := s.DB.CommentByNewsID(newsIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(comments)
	if err != nil {
		http.Error(w, "failed to marshal comments", http.StatusInternalServerError)
		return
	}

	w.Write(response)

}
