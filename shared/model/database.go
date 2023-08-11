package model

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
