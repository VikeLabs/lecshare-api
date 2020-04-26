package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"path"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const testingBucket = "assets-lecshare.oimo.ca"

func copyS3(srcKey string, srcBucket string, dstKey string, dstBucket string) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})
	if err != nil {
		log.Fatalln(err)
	}

	svc := s3.New(sess)

	_, err = svc.CopyObject(&s3.CopyObjectInput{
		Bucket:     aws.String(dstBucket),
		Key:        aws.String(dstKey),
		CopySource: aws.String(srcBucket + "/" + url.QueryEscape(srcKey)),
	})
	if err != nil {
		log.Fatalln(err)
	}

	err = svc.WaitUntilObjectExists(&s3.HeadObjectInput{
		Bucket: aws.String(dstBucket),
		Key:    aws.String(dstKey),
	})
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(">> Copied", dstKey, "to", dstBucket)
}

func deleteS3(key string, bucket string) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})
	if err != nil {
		log.Fatalln(err)
	}

	svc := s3.New(sess)

	_, err = svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Fatalln(err)
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(">> Deleted", key, "from", bucket)
}

func moveFile(key string, outKey string, s3object events.S3Entity) {
	copyS3(key, s3object.Bucket.Name, outKey, testingBucket)
	// Times out @ 6s, idk why. It has DeleteObject permissions
	//deleteS3(key, s3object.Bucket.Name)
}

func newFileHandler(ctx context.Context, event events.S3Event) error {
	for _, r := range event.Records {
		key, err := url.QueryUnescape(r.S3.Object.Key)
		if err != nil {
			return err
		}
		enc := base64.URLEncoding.WithPadding(base64.NoPadding)
		bytes, err := enc.DecodeString(strings.TrimSuffix(key, path.Ext(key)))
		if err != nil {
			return err
		}
		outKey := string(bytes) + ".json"
		fmt.Println(">> Storing", outKey)
		moveFile(key, outKey, r.S3)
	}

	return nil
}

func main() {
	lambda.Start(newFileHandler)
}
