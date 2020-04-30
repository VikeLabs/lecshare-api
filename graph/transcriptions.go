package graph

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/vikelabs/lecshare-api/graph/model"
)

func (r *Repository) GetTranscription(ctx context.Context, obj *model.Lecture) (*model.Transcription, error) {
	// when it retrieves the file, it stores it in memory (this buffer below) rather than on disk.
	buff := &aws.WriteAtBuffer{}

	downloader := s3manager.NewDownloader(r.Session)
	_, err := downloader.Download(buff,
		&s3.GetObjectInput{
			Bucket: r.AssetsBucketName,
			Key:    obj.Transcription,
		})

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// fmt.Println("Downloaded", *transcriptionFile, numBytes, "bytes")

	var transcriptionJSON model.TranscriptionJSON
	var transcription model.Transcription

	err = json.Unmarshal(buff.Bytes(), &transcriptionJSON)
	if err != nil {
		fmt.Println(err)
	}

	transcription.Transcripts = make([]*string, len(transcriptionJSON.Results.Transcripts))
	transcription.Words = make([]*model.TranscriptionWord, len(transcriptionJSON.Results.Items))

	for i, t := range transcriptionJSON.Results.Transcripts {
		transcription.Transcripts[i] = t.Transcript
	}
	for i, v := range transcriptionJSON.Results.Items {
		transcription.Words[i] = &model.TranscriptionWord{
			Type:      v.Type,
			Starttime: v.StartTime,
			Endtime:   v.EndTime,
			Word:      v.Alternatives[0].Content,
		}
	}

	return &transcription, nil
}
