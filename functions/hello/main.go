package main

import (
	"context"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
)

var h *httpadapter.HandlerAdapter

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	return h.ProxyWithContext(ctx, req)
}

func main() {
	h = httpadapter.New(playground.Handler("GraphQL playground", "/dev/query"))
	lambda.Start(handler)
}
