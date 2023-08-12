package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/dreamsofcode-io/reddit-api/shared/model"
)

type Database struct {
	client           *dynamodb.Client
	postTableName    string
	commentTableName string
	postIndexName    string
}

func NewDatabase(ctx context.Context) (*Database, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	postTableName := os.Getenv("POST_TABLE_NAME")
	commentTableName := os.Getenv("COMMENT_TABLE_NAME")
	postIndexName := os.Getenv("POST_INDEX_NAME")

	client := dynamodb.NewFromConfig(cfg)

	return &Database{
		client:           client,
		postTableName:    postTableName,
		commentTableName: commentTableName,
		postIndexName:    postIndexName,
	}, nil
}

func (d *Database) AddPostRecord(ctx context.Context, post model.PostRecord) error {
	item, err := attributevalue.MarshalMap(post)
	if err != nil {
		return fmt.Errorf("failed to marshal post: %w", err)
	}

	_, err = d.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(d.postTableName), Item: item,
	})
	if err != nil {
		return fmt.Errorf("failed to put item: %w", err)
	}

	return nil
}

func (d *Database) AddCommentRecord(ctx context.Context, comment model.CommentRecord) error {
	item, err := attributevalue.MarshalMap(comment)
	if err != nil {
		return fmt.Errorf("failed to marshal comment: %w", err)
	}

	_, err = d.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(d.commentTableName), Item: item,
	})
	if err != nil {
		return fmt.Errorf("failed to put item: %w", err)
	}

	return nil
}

func (d *Database) GetPostsForSubreddit(
	ctx context.Context,
	subreddit string,
	before *uint64,
) ([]model.PostRecord, error) {
	fmt.Println("subreddit: ", subreddit)
	keyEx := expression.Key("subreddit").Equal(expression.Value(subreddit))
	if before != nil {
		keyEx = keyEx.And(expression.Key("timestamp").LessThan(expression.Value(before)))
	}

	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		log.Printf("Couldn't build expression for query. Here's why: %v\n", err)
	}

	result, err := d.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(d.postTableName),
		IndexName:                 aws.String(d.postIndexName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		ScanIndexForward:          aws.Bool(false),
		Limit:                     aws.Int32(25),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query posts: %w", err)
	}

	var posts []model.PostRecord

	err = attributevalue.UnmarshalListOfMaps(result.Items, &posts)
	if err != nil {
		return nil, fmt.Errorf("failed to query posts: %w", err)
	}

	return posts, nil
}

func (d *Database) GetCommentsForPost(ctx context.Context, id string) (model.CommentRecord, error) {
	key, err := attributevalue.Marshal(id)
	if err != nil {
		return model.CommentRecord{}, fmt.Errorf("failed to marshal id: %w", err)
	}

	result, err := d.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(d.commentTableName),
		Key:       map[string]types.AttributeValue{"post_id": key},
	})
	if err != nil {
		return model.CommentRecord{}, fmt.Errorf("failed to get comment: %w", err)
	}

	var comment model.CommentRecord
	err = attributevalue.UnmarshalMap(result.Item, &comment)
	if err != nil {
		return model.CommentRecord{}, fmt.Errorf("failed to unmarshal comment: %w", err)
	}

	return comment, nil
}

func (d *Database) GetPost(ctx context.Context, id string) (model.PostRecord, error) {
	key, err := attributevalue.Marshal(id)
	if err != nil {
		return model.PostRecord{}, fmt.Errorf("failed to marshal id: %w", err)
	}

	result, err := d.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(d.postTableName),
		Key:       map[string]types.AttributeValue{"post_id": key},
	})
	if err != nil {
		return model.PostRecord{}, fmt.Errorf("failed to get post: %w", err)
	}

	var post model.PostRecord
	if err = attributevalue.UnmarshalMap(result.Item, &post); err != nil {
		return model.PostRecord{}, fmt.Errorf("failed to unmarshal post: %w", err)
	}

	return post, nil
}
