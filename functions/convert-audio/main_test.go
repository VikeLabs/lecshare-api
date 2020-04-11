package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

// Unit testing is on hold until I get it to work with pipes.

func TestProcessAudio(t *testing.T) {
	// k := "test.flac"
	k := "tell.flac"
	var s events.S3Entity
	s.Bucket.Name = "assets-lecshare.oimo.ca"
	processAudio(k, s)
}

// func TestDownloadS3Object(t *testing.T) {
// 	wg := sync.WaitGroup{}
// 	out, err := os.Create("test.flac")
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	go downloadS3("test.flac", "assets-lecshare.oimo.ca", out, &wg)
// 	wg.Wait()
// }

// func TestEncodeAudio(t *testing.T) {
// 	wg := sync.WaitGroup{}
// 	in, err := os.Open("test.flac")
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	out, err := os.Create("test-compressed.ogg")
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	go encodeAudio("test.flac", 128, in, out, &wg)
// 	wg.Wait()
// }

// func TestUploadS3Object(t *testing.T) {
// 	wg := sync.WaitGroup{}
// 	in, err := os.Open("test-compressed.ogg")
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	go uploadS3("test-compressed.ogg", "test.flac", 128, in, &wg)
// 	wg.Wait()
// }
