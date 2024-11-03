package main

import (
	"apigateway/APIGateWay/server"
	"apigateway/APIGateWay/server/middleware"
	"log"
	"net/http"
)

func main() {
	s := server.New()
	s.Register()

	log.Println("API Gateway service HTTP server is started on localhost:8083")

	http.ListenAndServe(":8083", middleware.RequestID(middleware.Logging(s.Router)))
}
