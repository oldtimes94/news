package main

import (
	"comments/middleware"
	"comments/server"
	"log"
	"net/http"
)

func main() {
	s, err := server.New("sqlite")
	if err != nil {
		log.Panicln(err)
	}
	s.Register()

	log.Println("Comments service HTTP server is started on localhost:8084")

	err = http.ListenAndServe(":8004", middleware.RequestID(middleware.Logging(s.Router)))
	if err != nil {
		log.Panicln(err)
	}

}
