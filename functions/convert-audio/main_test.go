package main

import (
	"os"
	"strconv"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/vikelabs/lecshare-api/utils"
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
	utils.DownloadS3(testFile+".flac", testingBucket, file)
	file.Close()
}

func TestProbeAudio(t *testing.T) {
	fileName := testFile + ".flac"
	file, err := os.Open(fileName)
	if err != nil {
		t.Error(err)
	}
	encoding, duration := utils.ProbeAudio(file, "")
	t.Log(encoding, "file is", strconv.Itoa(duration), "seconds.")
	file.Close()
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
	utils.EncodeAudio(128, "flac", "opus", in, out, "")
	in.Close()
	out.Close()
}

func TestUploadS3Object(t *testing.T) {
	key := testFile + "-compressed.ogg"
	file, err := os.Open(key)
	if err != nil {
		t.Error(err)
	}
	utils.UploadS3(key, testingBucket, "audio/ogg", file)
	file.Close()
}
