package main

import "testing"

func TestDownloadS3Object(t *testing.T) {
	downloadS3("test.flac", "assets-lecshare.oimo.ca")
}

func TestEncodeAudio(t *testing.T) {
	encodeAudio("test.flac", 128)
}

func TestUploadS3Object(t *testing.T) {
	uploadS3("test-compressed.ogg", "test.flac", 128)
}
