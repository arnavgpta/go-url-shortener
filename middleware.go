package main

import (
	"log"
	"net/http"
	"time"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("[%s] %s %s", r.Method, r.URL.Path, start.Format(time.RFC1123))
		next.ServeHTTP(w, r)
		log.Printf("Handled in %v\n", time.Since(start))
	})
}
