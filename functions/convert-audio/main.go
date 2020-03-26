// This is NOT a Lambda function yet, as there is no handler
// TODO Make this a proper Lambda function

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func downloadAWS(key string) {
	dir, _ := path.Split(key)

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		log.Fatalln(err)
	}
	file, err := os.Create(key)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	buff := &aws.WriteAtBuffer{}

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})

	downloader := s3manager.NewDownloader(sess)
	numBytes, err := downloader.Download(buff,
		&s3.GetObjectInput{
			Bucket: aws.String("assets-lecshare.oimo.ca"),
			Key:    aws.String(key),
		})

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Downloaded", key, numBytes, "bytes")
	_, err = file.Write(buff.Bytes())
	if err != nil {
		log.Fatalln(err)
	}
}

func uploadAWS(key string) {
	file, err := os.Open(key)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)

	uploader := s3manager.NewUploader(sess)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("assets-lecshare.oimo.ca"),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Uploaded", key)
}

func encodeAudio(filename string) string {
	baseName := strings.TrimSuffix(filename, path.Ext(filename))
	outName := baseName + "-compressed.ogg"

	out, err := exec.Command("/opt/bin/ffmpeg", "-i", filename, "-c:a", "libopus", "-ac", "1",
		"-b:a", "128k", outName).CombinedOutput()

	if err != nil {
		log.Fatalln(err, string(out))
	}
	fmt.Println("Created", outName)

	return outName
}

func main() {
	key := os.Args[1]

	// Where the magic happens
	downloadAWS(key)
	outKey := encodeAudio(key)
	uploadAWS(outKey)

	// Cleanup
	dir, _ := path.Split(key)
	err := os.RemoveAll(dir)
	if err != nil {
		log.Fatal(err)
	}
}
