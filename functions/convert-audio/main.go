// This is NOT a Lambda function yet, as there is no handler
// TODO Make this a proper Lambda function

package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func downloadS3(key string, bucket string) {
	dir, _ := path.Split(key)
	dir = "/tmp/" + dir

	if dir != "" {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			log.Fatalln(err)
		}
	}

	file, err := os.Create("/tmp/" + key)
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
			Bucket: aws.String(bucket),
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

func uploadS3(key string, oldKey string, bitrate int) {
	file, err := os.Open("/tmp/" + key)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)

	uploader := s3manager.NewUploader(sess)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String("assets-lecshare.oimo.ca"),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String("audio/ogg"),
		Metadata: aws.StringMap(map[string]string{
			"uncompressed-file-key": oldKey,
			"bitrate":               string(bitrate),
		}),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Uploaded", key)
}

// Bitrate is in kbps
func encodeAudio(filename string, bitrate int) string {
	baseName := strings.TrimSuffix(filename, path.Ext(filename))
	outName := baseName + "-compressed.ogg"

	cmd := exec.Command("/opt/ffmpeg/ffmpeg", "-y", "-i", "/tmp/"+filename, "-c:a", "libopus",
		"-ac", "1", "-b:a", strconv.Itoa(bitrate)+"k", "/tmp/"+outName)

	fmt.Println("Executing: " + cmd.Path + " " + strings.Join(cmd.Args, " "))

	out, err := cmd.CombinedOutput()

	if err != nil {
		log.Fatalln(err, string(out))
	}
	fmt.Println("Created", outName)

	return outName
}

func processAudio(key string, s3object events.S3Entity) {
	bitrate := 128

	// Where the magic happens
	downloadS3(key, s3object.Bucket.Name)
	outKey := encodeAudio(key, bitrate)
	uploadS3(outKey, key, bitrate)
	fmt.Print("\n")
}

func newAudioHandler(ctx context.Context, event events.S3Event) error {
	for _, r := range event.Records {
		key, err := url.QueryUnescape(r.S3.Object.Key)
		if err != nil {
			return err
		}
		fmt.Println("Processing ", key)
		processAudio(key, r.S3)
	}

	return nil
}

func main() {
	lambda.Start(newAudioHandler)
}
