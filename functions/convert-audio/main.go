package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"sync"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/transcribeservice"
)

var ffmpegDir string

const testingBucket = "assets-lecshare.oimo.ca"

// Only use this with downloader.Concurrency = 1, otherwise it will break.
type fakeWriterAt struct {
	w io.Writer
}

func (fw fakeWriterAt) WriteAt(p []byte, offset int64) (n int, err error) {
	// ignore 'offset' because we forced sequential downloads
	return fw.w.Write(p)
}

func downloadS3(key string, bucket string, outPipe *io.PipeWriter, wg *sync.WaitGroup) {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})

	downloader := s3manager.NewDownloader(sess)
	// Disable concurrency to sequentially stream the file
	downloader.Concurrency = 1
	numBytes, err := downloader.Download(fakeWriterAt{outPipe},
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(">> Downloaded", key, numBytes, "bytes")
	outPipe.Close()
	wg.Done()
}

func uploadS3(key string, oldKey string, bucket string, bitrate int, inPipe *io.PipeReader, wg *sync.WaitGroup) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)

	uploader := s3manager.NewUploader(sess)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        inPipe,
		ContentType: aws.String("audio/ogg"),
		Metadata: aws.StringMap(map[string]string{
			"uncompressed-file-key": oldKey,
			"bitrate":               string(bitrate),
		}),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(">> Uploaded", key)
	inPipe.Close()
	wg.Done()
}

// Bitrate is in kbps
func encodeAudio(bitrate int, inPipe *io.PipeReader, outPipe *io.PipeWriter, wg *sync.WaitGroup) {
	cmd := exec.Command(ffmpegDir+"ffmpeg", "-f", "flac", "-i", "pipe:", "-y", "-c:a", "libopus",
		"-ac", "1", "-b:a", strconv.Itoa(bitrate)+"k", "-f", "opus", "pipe:")

	fmt.Println(">> Executing: " + strings.Join(cmd.Args, " "))

	cmd.Stdin = inPipe
	cmd.Stdout = outPipe
	//cmd.Stderr = os.Stderr

	// Wait for command to complete.
	err := cmd.Run()

	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(">> Processed file.")
	inPipe.Close()
	outPipe.Close()
	wg.Done()
}

func processAudio(s3object events.S3Entity) {
	bitrate := 128
	key := s3object.Object.URLDecodedKey
	outKey := strings.TrimSuffix(key, path.Ext(key)) + ".ogg"

	inRead, inWrite := io.Pipe()
	outRead, outWrite := io.Pipe()
	wg := sync.WaitGroup{}

	// Where the magic happens
	wg.Add(3)
	go downloadS3(key, s3object.Bucket.Name, inWrite, &wg)
	go encodeAudio(bitrate, inRead, outWrite, &wg)
	go uploadS3(outKey, key, testingBucket, bitrate, outRead, &wg)
	fmt.Print("\n")
	wg.Wait()
}

func transcribeAudio(s3object events.S3Entity) {
	key := s3object.Object.URLDecodedKey
	// output file is jobName + ".json"
	jobName := strings.TrimSuffix(key, path.Ext(key))
	jobURI := "s3://" + s3object.Bucket.Name + "/" + s3object.Object.Key
	//jobUri := "https://uvic-transcribe-test1.s3-us-west-2.amazonaws.com/vikelabs_test1_15sec.flac"
	outBucket := testingBucket

	// open a new session
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})
	log.Println("Opening Transcribe session")
	transcriber := transcribeservice.New(sess)
	// exit if unable to create a Transcribe session
	if transcriber == nil {
		log.Fatalln("Unable to create Transcribe session")
	} else {
		log.Println("Transcribe session successfully created")
	}

	mediaformat := "flac"
	languagecode := "en-US"
	log.Println("Creating transcription job")
	var StrucMedia transcribeservice.Media
	StrucMedia.MediaFileUri = &jobURI
	_, err := transcriber.StartTranscriptionJob(&transcribeservice.StartTranscriptionJobInput{
		TranscriptionJobName: &jobName,
		Media:                &StrucMedia,
		MediaFormat:          &mediaformat,
		LanguageCode:         &languagecode,
		OutputBucketName:     &outBucket,
	})
	if err != nil {
		log.Fatalln("Got error building project: ", err)
	}

	log.Println("Successfully created transcription job for", key)
}

func newAudioHandler(ctx context.Context, event events.S3Event) error {
	for _, r := range event.Records {
		fmt.Println("Transcribing", r.S3.Object.URLDecodedKey)
		transcribeAudio(r.S3)
		fmt.Println("Processing ", r.S3.Object.URLDecodedKey)
		processAudio(r.S3)
	}

	return nil
}

func init() {
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		ffmpegDir = "/opt/ffmpeg/"
	}
}

func main() {
	lambda.Start(newAudioHandler)
}
