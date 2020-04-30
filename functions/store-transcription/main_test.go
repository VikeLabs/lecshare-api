package main

import (
	"encoding/base64"
	"testing"

	"github.com/vikelabs/lecshare-api/utils"
)

func TestCopyS3(t *testing.T) {
	enc := base64.URLEncoding.WithPadding(base64.NoPadding)

	outKey := "test.json"
	key := enc.EncodeToString([]byte("test")) + ".json"

	utils.CopyS3(key, transcriptionBucket, outKey, testingBucket)
}

func TestDeleteS3(t *testing.T) {
	enc := base64.URLEncoding.WithPadding(base64.NoPadding)

	key := enc.EncodeToString([]byte("test")) + ".json"

	utils.DeleteS3(key, transcriptionBucket)
}
