package main

import (
	"sync"
	"time"
)

var (
	urlStore     = make(map[string]urlData)
	reverseStore = make(map[string]string)
	mutex        = &sync.Mutex{}
)

type urlData struct {
	LongURL   string
	ExpiresAt time.Time
}
