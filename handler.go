package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		URL    string `json:"url"`
		Custom string `json:"custom"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existing URLMapping
	err := urlCollection.FindOne(ctx, bson.M{"long_url": req.URL}).Decode(&existing)
	if err == nil {
		resp := map[string]string{"short_url": "http://localhost:8080/" + existing.ShortCode}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}

	short := ""
	if req.Custom != "" {

		count, _ := urlCollection.CountDocuments(ctx, bson.M{"short_code": req.Custom})
		if count > 0 {
			http.Error(w, "Custom short code already exists", http.StatusConflict)
			return
		}
		short = req.Custom
	} else {
		short = generateShortKey()
	}

	entry := URLMapping{
		ShortCode: short,
		LongURL:   req.URL,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	_, err = urlCollection.InsertOne(ctx, entry)
	if err != nil {
		log.Println("MongoDB insert error:", err)
		http.Error(w, "Failed to save URL", http.StatusInternalServerError)
		return
	}

	resp := map[string]string{"short_url": "http://localhost:8080/" + short}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	short := r.URL.Path[1:]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result URLMapping
	err := urlCollection.FindOne(ctx, bson.M{"short_code": short}).Decode(&result)
	if err != nil {
		http.Error(w, "URL not found or expired", http.StatusNotFound)
		return
	}

	if result.ExpiresAt.Before(time.Now()) {
		http.Error(w, "URL has expired", http.StatusGone)
		return
	}

	http.Redirect(w, r, result.LongURL, http.StatusFound)
}
