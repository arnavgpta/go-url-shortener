package main

import (
	"log"
	"net/http"
)

func main() {
	initMongo()

	mux := http.NewServeMux()
	mux.HandleFunc("/shorten", shortenHandler)
	mux.HandleFunc("/", redirectHandler)

	log.Println("Server is running on http://localhost:8080")
	err := http.ListenAndServe(":8080", loggingMiddleware(mux))
	if err != nil {
		log.Fatal(err)
	}
}
