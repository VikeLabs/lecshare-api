package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/vikelabs/lecshare-api/graph/generated"
	"github.com/vikelabs/lecshare-api/graph/model"
)

func (r *lectureResolver) Transcription(ctx context.Context, obj *model.Lecture) (*model.Transcription, error) {
	// when it retrieves the file, it stores it in memory (this buffer below) rather than on disk.
	buff := &aws.WriteAtBuffer{}

	downloader := s3manager.NewDownloader(r.Session)
	_, err := downloader.Download(buff,
		&s3.GetObjectInput{
			Bucket: r.BucketName,
			Key:    obj.Transcription,
		})

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// fmt.Println("Downloaded", *transcriptionFile, numBytes, "bytes")

	var transcription model.Transcription

	err = json.Unmarshal(buff.Bytes(), &transcription.Sections)
	if err != nil {
		fmt.Println(err)
	}

	return &transcription, nil
}

// Class returns generated.ClassResolver implementation.
func (r *Resolver) Class() generated.ClassResolver { return &classResolver{r} }

// Course returns generated.CourseResolver implementation.
func (r *Resolver) Course() generated.CourseResolver { return &courseResolver{r} }

// Lecture returns generated.LectureResolver implementation.
func (r *Resolver) Lecture() generated.LectureResolver { return &lectureResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// School returns generated.SchoolResolver implementation.
func (r *Resolver) School() generated.SchoolResolver { return &schoolResolver{r} }

type classResolver struct{ *Resolver }
type courseResolver struct{ *Resolver }
type lectureResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type schoolResolver struct{ *Resolver }
