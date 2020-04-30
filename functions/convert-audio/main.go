package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/transcribeservice"
	"github.com/vikelabs/lecshare-api/utils"
)

var ffmpegDir string
var tmpDir string

const testingBucket = "assets-lecshare.oimo.ca"
const transcriptionBucket = "lecshare-transcriptions"

func processAudio(key string, s3object events.S3Entity) {
	keyPath := path.Dir(key)
	//keyFile := path.Base(key)

	if keyPath != "" {
		err := os.MkdirAll(tmpDir+keyPath, 0755)
		if err != nil {
			log.Fatalln(err)
		}
	}

	inFile, err := os.Create(tmpDir + key)
	if err != nil {
		log.Fatalln(err)
	}

	utils.DownloadS3(key, s3object.Bucket.Name, inFile)

	inFile.Close()
	inFile, err = os.Open(tmpDir + key)
	if err != nil {
		log.Fatalln(err)
	}

	bitrate := 128
	inCodec, duration := utils.ProbeAudio(inFile, ffmpegDir)
	fmt.Println(">> File is type", inCodec, "and length", duration, "seconds.")

	//  Codec     ext
	finalCodecs := map[string]string{
		"opus": ".ogg",
		"mp3":  ".mp3",
	}

	for outCodec, extension := range finalCodecs {
		outKey := strings.TrimSuffix(key, path.Ext(key)) + extension
		outFile, err := os.Create(tmpDir + outKey)
		if err != nil {
			log.Fatalln(err)
		}
		inFile.Seek(0, os.SEEK_SET)

		utils.EncodeAudio(bitrate, inCodec, outCodec, inFile, outFile, ffmpegDir)
		outFile.Close()
		outFile, err = os.Open(tmpDir + outKey)
		if err != nil {
			log.Fatalln(err)
		}

		mime := "audio/" + strings.TrimLeft(extension, ".")
		utils.UploadS3(outKey, testingBucket, mime, outFile)
	}

	fmt.Println("Transcribing", key)
	transcribeAudio(key, s3object, inCodec)
}

func transcribeAudio(key string, s3object events.S3Entity, codec string) {
	// encode original key in jobName so we can smuggle it into the output
	// no padding because jobName allows a limited char set: ^[0-9a-zA-Z._-]+ up to 200 chars
	enc := base64.URLEncoding.WithPadding(base64.NoPadding)

	jobName := enc.EncodeToString([]byte(strings.TrimSuffix(key, path.Ext(key))))
	jobURI := "s3://" + s3object.Bucket.Name + "/" + s3object.Object.Key
	outBucket := transcriptionBucket

	// open a new session
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})
	fmt.Println(">> Opening Transcribe session for", key)
	transcriber := transcribeservice.New(sess)
	// exit if unable to create a Transcribe session
	if transcriber == nil {
		log.Fatalln("Unable to create Transcribe session")
	} else {
		fmt.Println(">> Transcribe session successfully created")
	}

	mediaformat := "flac"
	languagecode := "en-US"
	var StrucMedia transcribeservice.Media
	StrucMedia.MediaFileUri = &jobURI

	fmt.Println(">> Creating transcription job")

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

	fmt.Println("Successfully created transcription job for", key)
}

func newAudioHandler(ctx context.Context, event events.S3Event) error {
	for _, r := range event.Records {
		key, err := url.QueryUnescape(r.S3.Object.Key)
		if err != nil {
			return err
		}
		fmt.Println("Processing", key)
		processAudio(key, r.S3)
	}

	return nil
}

func init() {
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		ffmpegDir = "/opt/ffmpeg/"
		tmpDir = "/tmp/"
	}
}

func main() {
	lambda.Start(newAudioHandler)
}
