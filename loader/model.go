package main

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

type Comment struct {
	ID        string    `json:"id"`
	Author    string    `json:"author"`
	Timestamp time.Time `json:"time"`
	Content   string    `json:"comment"`
	Points    int       `json:"points"`
	IsDeleted bool      `json:"isDeleted"`
	Children  []Comment `json:"children"`
}

type PostRecord struct {
	ID        string `dynamodbav:"post_id"`
	Subreddit string `dynamodbav:"subreddit"`
	Timestamp uint64 `dynamodbav:"timestamp"`
	Data      string `dynamodbav:"data"`
}

type CommentRecord struct {
	ID   string `dynamodbav:"post_id"`
	Data string `dynamodbav:"data"`
}
