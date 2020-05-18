package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/go-playground/validator/v10"
	"github.com/guregu/dynamo"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jws"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/vikelabs/lecshare-api/graph"
	"github.com/vikelabs/lecshare-api/graph/generated"
	"github.com/vikelabs/lecshare-api/utils/bunnycdn"
)

var h *httpadapter.HandlerAdapter
var keys *jwk.Set

func lambdaHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// https://docs.aws.amazon.com/cognito/latest/developerguide/amazon-cognito-user-pools-using-tokens-verifying-a-jwt.html
	authHeader := req.Headers["Authorization"]
	authHeader = strings.TrimSpace(authHeader)
	authHeaderSegments := strings.Split(authHeader, " ")
	if len(authHeaderSegments) == 2 && authHeaderSegments[0] == "Bearer" {
		payloadBytes := []byte(authHeaderSegments[1])
		_, err := jws.VerifyWithJWKSet(payloadBytes, keys, jws.DefaultJWKAcceptor)
		if err == nil {
			fmt.Println("Validily signed access token")
			token, err := jwt.ParseBytes(payloadBytes,
				jwt.WithAudience("2rt4ddvsjbndmovp9ug6tur2m"),
				jwt.WithIssuer("https://cognito-idp.us-west-2.amazonaws.com/us-west-2_LTj42jgMX"),
			)
			if err != nil {
				fmt.Println(err)
			} else {
				// inject the token related info into context.
				ctx = context.WithValue(ctx, "jwt", token.AsMap)
			}
		} else {
			fmt.Println("Invalid access_token.")
		}
	}
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
	cdn := os.Getenv("CDN")
	tableName := os.Getenv("TABLE_NAME")

	validate := validator.New()
	presigner := bunnycdn.Generator{
		APIKey:   os.Getenv("CDN_API_KEY"),
		Hostname: cdn,
	}

	k, err := jwk.Fetch("https://cognito-idp.us-west-2.amazonaws.com/us-west-2_LTj42jgMX/.well-known/jwks.json")
	if err != nil {
		log.Printf("failed to parse JWK: %s", err)
		return
	}
	keys = k

	srv := handler.New(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		Repository: graph.Repository{
			DynamoDB:              db,
			TableName:             &tableName,
			Session:               session,
			AssetsBucketName:      &bucketName,
			ProcessingBucketName:  &processingBucketName,
			Validate:              validate,
			CDN:                   &cdn,
			PresignedURLGenerator: &presigner,
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

	h = httpadapter.New(srv)

	lambda.Start(lambdaHandler)
}
