package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func addPost(ctx context.Context, post Post, db *Database) error {
	comments := post.Comments
	post.Comments = nil

	postData, err := json.Marshal(post)
	if err != nil {
		return fmt.Errorf("failed to marshal post: %w", err)
	}

	commentData, err := json.Marshal(comments)
	if err != nil {
		return fmt.Errorf("failed to marshal comments: %w", err)
	}

	postRecord := PostRecord{
		ID:        post.ID,
		Subreddit: post.Subreddit,
		Timestamp: post.TimestampMillis / 1000,
		Data:      string(postData),
	}

	commentRecord := CommentRecord{
		ID:   post.ID,
		Data: string(commentData),
	}

	if err := db.AddCommentRecord(ctx, commentRecord); err != nil {
		return fmt.Errorf("failed to add comment record: %w", err)
	}

	if err := db.AddPostRecord(ctx, postRecord); err != nil {
		return fmt.Errorf("failed to add post record: %w", err)
	}

	return nil
}

func handle(ctx context.Context, event events.SQSEvent, db *Database) error {
	for _, record := range event.Records {
		var post Post
		if err := json.Unmarshal([]byte(record.Body), &post); err != nil {
			fmt.Println("Failed to unmarshal post: ", record.Body)
			return fmt.Errorf("failed to unmarshal post: %w", err)
		}

		if err := addPost(ctx, post, db); err != nil {
			return fmt.Errorf("failed to add post: %w", err)
		}
	}

	return nil
}

func main() {
	timeout := time.Duration(3) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	db, err := NewDatabase(ctx)
	if err != nil {
		panic(err)
	}

	lambda.Start(func(ctx context.Context, event events.SQSEvent) error {
		return handle(ctx, event, db)
	})
}
