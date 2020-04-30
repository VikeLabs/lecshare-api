package utils

import (
	"fmt"
	"io"
	"log"
	"net/url"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

//TODO better error handling, we just log.Fatalln(err) right now.

// DownloadS3 downloads key from bucket and writes it to output
func DownloadS3(key string, bucket string, output io.Writer) {
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

	output.Write(buff.Bytes())

	fmt.Println(">> Downloaded", key, numBytes, "bytes")
}

//UploadS3 writes input to key in bucket with content-type mime.
func UploadS3(key string, bucket string, mime string, input io.Reader) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)

	uploader := s3manager.NewUploader(sess)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        input,
		ContentType: aws.String(mime),
		// TODO Make metadata an argument of type map[string]string
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

// CopyS3 copies srcKey in srcBucket to dstKey in dstBucket
func CopyS3(srcKey string, srcBucket string, dstKey string, dstBucket string) {
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

// DeleteS3 deletes key from bucket
// I'm not sure if this works on Lambda, sometimes it times out on WaitUntilObjectNotExists
func DeleteS3(key string, bucket string) {
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
