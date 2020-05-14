package main

import (
	"context"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/go-playground/validator/v10"
	"github.com/guregu/dynamo"
	"github.com/vikelabs/lecshare-api/graph"
	"github.com/vikelabs/lecshare-api/graph/generated"
)

var h *httpadapter.HandlerAdapter

func lambdaHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// passes in the APIGatewayProxyRequest as a context.
	// c, _ := core.GetAPIGatewayContextFromContext(ctx)
	res, err := h.ProxyWithContext(ctx, req)
	res.Headers["Access-Control-Allow-Origin"] = "*"
	res.Headers["Access-Control-Allow-Headers"] = "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token,X-Amz-User-Agent"
	return res, err
}

func main() {
	session := session.New(&aws.Config{Region: aws.String("us-west-2")})
	db := dynamo.New(session)

	bucketName := os.Getenv("BUCKET_NAME")
	processingBucketName := os.Getenv("PROCESSING_BUCKET_NAME")
	// cdn := os.Getenv("CDN")
	tableName := os.Getenv("TABLE_NAME")

	validate := validator.New()

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
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

	h = httpadapter.New(srv)

	lambda.Start(lambdaHandler)
}
