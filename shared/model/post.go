package model

import "time"

type Post struct {
	ID              string    `json:"id"`
	Subreddit       string    `json:"subreddit"`
	DataType        string    `json:"dataType"`
	DataURL         string    `json:"dataURL"`
	IsPromoted      bool      `json:"isPromoted"`
	IsGallery       bool      `json:"isGallery"`
	Title           string    `json:"title"`
	Timestamp       time.Time `json:"timestamp"`
	TimestampMillis uint64    `json:"timestamp_millis"`
	Author          string    `json:"author"`
	URL             string    `json:"url"`
	Points          int       `json:"points"`
	ScrapedAt       time.Time `json:"scrapedAt"`
	Comments        []Comment `json:"comments,omitempty"`
}
