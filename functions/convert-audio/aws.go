package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func downloadS3(key string, bucket string, file *os.File) {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})

	buff := aws.WriteAtBuffer{}

	downloader := s3manager.NewDownloader(sess)
	numBytes, err := downloader.Download(&buff,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
	if err != nil {
		log.Fatalln(err)
	}

	file.Write(buff.Bytes())

	fmt.Println(">> Downloaded", key, numBytes, "bytes")
	file.Close()
}

func uploadS3(key string, bucket string, mime string, reader io.Reader) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)

	uploader := s3manager.NewUploader(sess)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        reader,
		ContentType: aws.String("audio/" + mime),
		// Metadata: aws.StringMap(map[string]string{
		// 	"uncompressed-file-key": oldKey,
		// 	"bitrate":               string(bitrate),
		// }),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(">> Uploaded", key)
}
