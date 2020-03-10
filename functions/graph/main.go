package main

import (
	"context"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/vikelabs/lecshare-api/graph"
	"github.com/vikelabs/lecshare-api/graph/generated"
)

var h *httpadapter.HandlerAdapter

func lambdaHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// passes in the APIGatewayProxyRequest as a context.
	// c, _ := core.GetAPIGatewayContextFromContext(ctx)
	return h.ProxyWithContext(ctx, req)
}

func main() {
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	h = httpadapter.New(srv)

	lambda.Start(lambdaHandler)
}
