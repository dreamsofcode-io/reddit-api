package main

import (
	"context"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/dreamsofcode-io/reddit-api/shared/database"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	db, err := database.NewDatabase(context.Background())
	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()

	router.HandleFunc("/", helloWorld()).Methods(http.MethodGet)
	router.HandleFunc("/r/{subreddit}", getSubreddit(db)).Methods(http.MethodGet)
	router.HandleFunc("/posts/{id}", getPost(db)).Methods(http.MethodGet)

	app := httpadapter.New(handlers.LoggingHandler(os.Stdout, router))

	lambda.Start(app.ProxyWithContext)
}
