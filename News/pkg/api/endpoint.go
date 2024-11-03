package api

import (
	"GoNews/pkg/sqlbuild"
	"GoNews/pkg/storage"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (api *API) postsHandler(w http.ResponseWriter, r *http.Request) {

	countQuery := sqlbuild.NewsQuery(r.URL, true, 0)
	count, err := api.db.Count(countQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	page := 1
	pageFromQuery := r.URL.Query().Get("page")
	if pageFromQuery != "" {
		page, err = strconv.Atoi(pageFromQuery)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	totalPages := (count + sqlbuild.ItemPerPage - 1) / sqlbuild.ItemPerPage
	offset := (page - 1) * sqlbuild.ItemPerPage
	pagen := storage.NewPagination(totalPages, page)

	query := sqlbuild.NewsQuery(r.URL, false, offset)
	posts := make([]storage.NewsPost, 0, 10)

	posts, err = api.db.Posts(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)

	pagenResponse, err := json.Marshal(pagen)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(pagenResponse)

}

func (api *API) DetailPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	post, err := api.db.Post(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(response)

}
