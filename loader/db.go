package main

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Database struct {
	client           *dynamodb.Client
	postTableName    string
	commentTableName string
}

func NewDatabase(ctx context.Context) (*Database, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	postTableName := os.Getenv("POST_TABLE_NAME")
	commentTableName := os.Getenv("COMMENT_TABLE_NAME")

	client := dynamodb.NewFromConfig(cfg)

	return &Database{
		client:           client,
		postTableName:    postTableName,
		commentTableName: commentTableName,
	}, nil
}

func (d *Database) AddPostRecord(ctx context.Context, post PostRecord) error {
	item, err := attributevalue.MarshalMap(post)
	if err != nil {
		return err
	}

	_, err = d.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(d.postTableName), Item: item,
	})

	return nil
}

func (d *Database) AddCommentRecord(ctx context.Context, comment CommentRecord) error {
	item, err := attributevalue.MarshalMap(comment)
	if err != nil {
		return err
	}

	_, err = d.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(d.commentTableName), Item: item,
	})

	return nil
}
