package main

import (
	"commentvalidator/pkg/server"
	"log"
	"net/http"
)

func main() {

	s, err := server.New()
	if err != nil {
		log.Panicln(err)
	}
	s.Register()

	log.Println("CommentValidator service HTTP server is started on localhost:8003")

	err = http.ListenAndServe(":8003", server.Logging(s.Router))
	if err != nil {
		log.Panicln(err)
	}
}
