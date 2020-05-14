package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-playground/validator/v10"
	"github.com/guregu/dynamo"
	"github.com/joho/godotenv"
	"github.com/vikelabs/lecshare-api/graph"
	"github.com/vikelabs/lecshare-api/graph/generated"
	"github.com/vikelabs/lecshare-api/utils"
)

var (
	port int
	host string
)

func main() {
	flag.IntVar(&port, "p", 8080, "specify port to use")
	flag.StringVar(&host, "h", "localhost", "specfiy host to bind to")

	flag.Parse()
	session := session.New(&aws.Config{Region: aws.String("us-west-2")})
	db := dynamo.New(session)

	// for local development only
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Preparation for resolvers.
	bucketName := os.Getenv("BUCKET_NAME")
	processingBucketName := os.Getenv("PROCESSING_BUCKET_NAME")
	cdn := os.Getenv("CDN")
	tableName := os.Getenv("TABLE_NAME")

	validate := validator.New()

	// Initialize GraphQL resolvers
	srv := handler.New(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		Repository: graph.Repository{
			DynamoDB:             db,
			TableName:            &tableName,
			Session:              session,
			AssetsBucketName:     &bucketName,
			ProcessingBucketName: &processingBucketName,
			Validate:             validate,
		},
	}}))

	var mb int64 = 1 << 20

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{
		MaxMemory:     128 * mb,
		MaxUploadSize: 128 * mb,
	})

	srv.SetQueryCache(lru.New(1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})

	// define routes for development.
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	// host on user defined port / default port.
	log.Printf("connect to http://%s:%d/ for GraphQL playground", host, port)
	log.Printf("connect to http://%s:%d/query for GraphQL endpoint", host, port)
	log.Fatal(http.ListenAndServe(host+":"+strconv.Itoa(port), nil))
}
