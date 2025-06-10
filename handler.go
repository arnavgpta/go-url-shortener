package main

import (
	"encoding/json"
	"net/http"
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

	short := ""

	mutex.Lock()
	defer mutex.Unlock()

	if existingShort, exists := reverseStore[req.URL]; exists {
		resp := map[string]string{"short_url": "http://localhost:8080/" + existingShort}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}

	if req.Custom != "" {
		if _, exists := urlStore[req.Custom]; exists {
			http.Error(w, "Custom short code already exists", http.StatusConflict)
			return
		}
		short = req.Custom
	} else {
		short = generateShortKey()
	}

	urlStore[short] = req.URL
	reverseStore[req.URL] = short

	resp := map[string]string{"short_url": "http://localhost:8080/" + short}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	short := r.URL.Path[1:]

	mutex.Lock()
	longURL, ok := urlStore[short]
	mutex.Unlock()

	if !ok {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, longURL, http.StatusFound)
}
