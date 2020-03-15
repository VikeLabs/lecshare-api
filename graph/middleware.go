package graph

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Middleware is an example usage of a simple middlware.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "Message", "Hello World")
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// CorsMiddleware is a help function get around CORS
// TODO do this better with authentication
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

// GetResource gets an expiring url for any file in the S3 bucket
// TODO limit filesystem access for security purposes
func GetResource(filename string, expire time.Duration) (string, error) {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})

	svc := s3.New(sess)

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String("assets-lecshare.oimo.ca"),
		Key:    aws.String(filename),
	})

	// Gets presigned URL with arbitrary lifetime
	urlStr, err := req.Presign(expire * time.Minute)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return urlStr, nil
}
