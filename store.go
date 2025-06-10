package main

import "sync"

var (
	urlStore = make(map[string]string)
	mutex    = &sync.Mutex{}
)
