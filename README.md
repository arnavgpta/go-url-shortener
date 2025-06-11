# Go URL Shortener 

A lightweight and scalable URL shortener built with Go and MongoDB.

## Features

- Generate short links from long URLs
- Support for custom short codes (optional)
- Automatic expiry for short URLs
- MongoDB-based persistence
- Middleware-based request logging
- Concurrency-safe design using Go's built-in features

## Tech Stack

- Language: Go
- Database: MongoDB
- Tools: Postman, Git

---

## API Endpoints

### 1. POST `/shorten`

**Request:**

```json
{
  "url": "https://example.com/long-url",
  "custom": "mycustom"   // optional
}
