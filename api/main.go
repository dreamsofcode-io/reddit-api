package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/dreamsofcode-io/reddit-api/shared/database"
	"github.com/dreamsofcode-io/reddit-api/shared/model"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	db, err := database.NewDatabase(context.Background())
	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World 👋!"))
	}).Methods(http.MethodGet)

	router.HandleFunc("/r/{subreddit}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		subreddit := fmt.Sprintf("r/%s", vars["subreddit"])

		res, err := db.GetPostsForSubreddit(r.Context(), subreddit)
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
	})

	router.HandleFunc("/post/{id}", func(w http.ResponseWriter, r *http.Request) {
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

		if err = json.NewEncoder(w).Encode(post); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println("failed to encode post: ", err)
			return
		}
	}).Methods(http.MethodGet)

	router.HandleFunc("/post/{id}/comments", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		postID := vars["id"]

		res, err := db.GetCommentsForPost(r.Context(), postID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println("failed to get comments for post: ", err)
			return
		}

		var comments []model.Comment
		if err = json.Unmarshal([]byte(res.Data), &comments); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println("failed to unmarshal comment: ", err)
			return
		}

		if err = json.NewEncoder(w).Encode(comments); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println("failed to encode comments: ", err)
			return
		}

	}).Methods(http.MethodGet)

	app := httpadapter.New(handlers.LoggingHandler(os.Stdout, router))

	lambda.Start(app.ProxyWithContext)
}
