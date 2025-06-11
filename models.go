package main

import "time"

type URLMapping struct {
	ShortCode string    `bson:"short_code"`
	LongURL   string    `bson:"long_url"`
	CreatedAt time.Time `bson:"created_at"`
	ExpiresAt time.Time `bson:"expires_at"`
}
