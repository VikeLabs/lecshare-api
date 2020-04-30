package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/vikelabs/lecshare-api/utils"
)

const testingBucket = "assets-lecshare.oimo.ca"
const transcriptionBucket = "lecshare-transcriptions"

func moveFile(key string, outKey string, s3object events.S3Entity) {
	utils.CopyS3(key, s3object.Bucket.Name, outKey, testingBucket)
	// Times out, idk why. It has DeleteObject permissions
	//utils.deleteS3(key, s3object.Bucket.Name)
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
