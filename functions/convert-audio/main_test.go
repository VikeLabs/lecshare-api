package main

import (
	"os"
	"strconv"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

const testFile = "test_file"

// The big test
func TestProcessAudio(t *testing.T) {
	k := testFile + ".flac"
	var s events.S3Entity
	s.Bucket.Name = testingBucket
	s.Object.Key = k
	processAudio(k, s)
}

func TestTranscribeAudio(t *testing.T) {
	// Will fail if you use the same job name repeatedly
	k := testFile + ".flac"
	var s events.S3Entity
	s.Bucket.Name = testingBucket
	s.Object.Key = k
	transcribeAudio(k, s, "flac")
}

func TestDownloadS3Object(t *testing.T) {
	file, err := os.Create(testFile + ".flac")
	if err != nil {
		t.Error(err)
	}
	downloadS3(testFile+".flac", testingBucket, file)
}

func TestProbeAudio(t *testing.T) {
	fileName := testFile + ".flac"
	reader, err := os.Open(fileName)
	if err != nil {
		t.Error(err)
	}
	encoding, duration := probeAudio(reader)
	t.Log(encoding, "file is", strconv.Itoa(duration), "seconds.")
}

func TestEncodeAudio(t *testing.T) {
	in, err := os.Open(testFile + ".flac")
	if err != nil {
		t.Error(err)
	}
	out, err := os.Create(testFile + "-compressed.ogg")
	if err != nil {
		t.Error(err)
	}
	encodeAudio(128, "flac", "opus", in, out)
}

func TestUploadS3Object(t *testing.T) {
	key := testFile + "-compressed.ogg"
	in, err := os.Open(key)
	if err != nil {
		t.Error(err)
	}
	uploadS3(key, testingBucket, "audio/ogg", in)
}
