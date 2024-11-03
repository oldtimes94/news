package server

import (
	"io"
	"log"
	"net/http"
	"strings"
)

func (s *Server) news(w http.ResponseWriter, r *http.Request) {

	req, err := http.NewRequest("GET", News.Path(r.RequestURI), nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	req.Header.Set("X-Request-ID", r.Context().Value("request_id").(string))
	client := &http.Client{}

	response, err := client.Do(req)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(body)

}

func (s *Server) comment(w http.ResponseWriter, r *http.Request) {

	params := strings.Join([]string{"/validate?", r.URL.RawQuery}, "")
	req, err := http.NewRequest("GET", CommentValidator.Path(params), nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	req.Header.Set("X-Request-ID", r.Context().Value("request_id").(string))
	client := &http.Client{}

	validateResponse, err := client.Do(req)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if validateResponse.StatusCode != http.StatusOK {
		http.Error(w, validateResponse.Status, validateResponse.StatusCode)
		return
	}

	params = strings.Join([]string{"/news?", r.URL.RawQuery}, "")
	req, err = http.NewRequest("POST", Comments.Path(params), nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	req.Header.Set("X-Request-ID", r.Context().Value("request_id").(string))
	client = &http.Client{}
	postCommentResponse, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer postCommentResponse.Body.Close()

	if postCommentResponse.StatusCode != http.StatusOK {
		http.Error(w, postCommentResponse.Status, postCommentResponse.StatusCode)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) newsDetail(w http.ResponseWriter, r *http.Request) {

	resp := make(chan []byte)
	errors := make(chan error, 2)

	go func() {

		req, err := http.NewRequest("GET", News.Path(r.RequestURI), nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		req.Header.Set("X-Request-ID", r.Context().Value("request_id").(string))
		client := &http.Client{}

		newsResponse, err := client.Do(req)
		if err != nil {
			errors <- err
			return
		}
		defer newsResponse.Body.Close()

		body, err := io.ReadAll(newsResponse.Body)
		if err != nil {
			errors <- err
			return
		}

		resp <- body
	}()

	go func() {

		req, err := http.NewRequest("GET", Comments.Path(r.RequestURI), nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		req.Header.Set("X-Request-ID", r.Context().Value("request_id").(string))
		client := &http.Client{}

		commentsResponse, err := client.Do(req)
		if err != nil {
			errors <- err
			return
		}
		defer commentsResponse.Body.Close()

		body, err := io.ReadAll(commentsResponse.Body)
		if err != nil {
			errors <- err
			return
		}

		resp <- body
	}()

	var result []byte
	for i := 0; i < 2; i++ {
		select {
		case res := <-resp:
			result = append(result, res...)
		case err := <-errors:
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}
