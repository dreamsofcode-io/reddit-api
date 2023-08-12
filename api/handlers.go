package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dreamsofcode-io/reddit-api/shared/database"
	"github.com/dreamsofcode-io/reddit-api/shared/model"
	"github.com/gorilla/mux"
)

func helloWorld() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World 👋!"))
	}
}

func getAfterTimestamp(ctx context.Context, db *database.Database, id string) (*uint64, error) {
	if id == "" {
		return nil, nil
	}

	row, err := db.GetPost(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	return &row.Timestamp, nil
}

func getSubreddit(db *database.Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		subreddit := fmt.Sprintf("r/%s", vars["subreddit"])

		before, err := getAfterTimestamp(r.Context(), db, r.URL.Query().Get("after"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println("failed to get after timestamp: ", err)
			return
		}

		res, err := db.GetPostsForSubreddit(r.Context(), subreddit, before)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println("failed to get posts for subreddit: ", err)
			return
		}

		posts := make([]model.Post, len(res))
		for i, x := range res {
			var post model.Post
			if err = json.Unmarshal([]byte(x.Data), &post); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Println("failed to unmarshal post: ", err)
				return
			}

			posts[i] = post
		}

		if err = json.NewEncoder(w).Encode(posts); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println("failed to encode posts: ", err)
			return
		}
	}
}

func getPost(db *database.Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		postID := vars["id"]

		res, err := db.GetPost(r.Context(), postID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println("failed to get post: ", err)
			return
		}

		var post model.Post
		if err = json.Unmarshal([]byte(res.Data), &post); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println("failed to unmarshal post: ", err)
			return
		}

		commentRow, err := db.GetCommentsForPost(r.Context(), postID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println("failed to get comments for post: ", err)
			return
		}

		if err = json.Unmarshal([]byte(commentRow.Data), &post.Comments); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println("failed to unmarshal comment: ", err)
			return
		}

		if err = json.NewEncoder(w).Encode(post); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println("failed to encode post: ", err)
			return
		}
	}
}
