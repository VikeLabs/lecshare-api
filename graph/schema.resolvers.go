// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
package graph

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/vikelabs/lecshare-api/graph/generated"
	"github.com/vikelabs/lecshare-api/graph/model"
)

func (r *classResolver) Lectures(ctx context.Context, obj *model.Class) ([]*model.Lecture, error) {
	return []*model.Lecture{
		&model.Lecture{
			Name:     "Introduction",
			Datetime: "Feb, 12, 2020",
			Duration: 3600,
		},
		&model.Lecture{
			Name:     "Final",
			Datetime: "Feb 13, 2020",
			Duration: 3600,
		},
	}, nil
}

func (r *lectureResolver) Transcription(ctx context.Context, obj *model.Lecture) (*model.Transcription, error) {
	buff := &aws.WriteAtBuffer{}

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})
	const transcriptionFile = "vikelabs_test1.json"

	downloader := s3manager.NewDownloader(sess)
	_, err := downloader.Download(buff,
		&s3.GetObjectInput{
			Bucket: aws.String("assets-lecshare.oimo.ca"),
			Key:    aws.String(transcriptionFile),
		})

	// fmt.Println("Downloaded", transcriptionFile, numBytes, "bytes")

	var transcription model.Transcription

	err = json.Unmarshal(buff.Bytes(), &transcription.Sections)
	if err != nil {
		fmt.Println(err)
	}

	return &transcription, nil
}

func (r *queryResolver) Schools(ctx context.Context, shortname *string) ([]*model.School, error) {
	fmt.Println(ctx.Value("Auth"))

	// TODO retrieve list of schools on Lecshare.
	if shortname != nil {
		if *shortname == "UVIC" {
			return []*model.School{
				&model.School{
					Name:      "University of Victoria",
					Shortname: "UVIC",
				},
			}, nil
		} else if *shortname == "VLABS" {
			return []*model.School{
				&model.School{
					Name:      "VikeLabs",
					Shortname: "VLABS",
				},
			}, nil
		}
	}

	return []*model.School{
		&model.School{
			Name:      "University of Victoria",
			Shortname: "UVIC",
		},
		&model.School{
			Name:      "VikeLabs",
			Shortname: "VLABS",
		},
	}, nil

}

func (r *schoolResolver) Classes(ctx context.Context, obj *model.School, typeArg *string) ([]*model.Class, error) {
	// TODO retrive list of classes available in the school.
	if obj.Shortname == "UVIC" {
		return []*model.Class{
			&model.Class{
				Title: "Foundations of Programming II",
				Code:  "CSC 115",
				Instructor: &model.User{
					FirstName: "Bill",
					LastName:  "Bird",
					Prefix:    "Dr",
					Role:      "Instructor",
				},
			},
		}, nil
	}
	return []*model.Class{
		&model.Class{
			Title: "Introduction to Git",
			Code:  "GIT 101",
			Instructor: &model.User{
				FirstName: "Aomi",
				LastName:  "Jokoji",
				Prefix:    "",
				Role:      "Student",
			},
		},
	}, nil
}

func (r *Resolver) Class() generated.ClassResolver     { return &classResolver{r} }
func (r *Resolver) Lecture() generated.LectureResolver { return &lectureResolver{r} }
func (r *Resolver) Query() generated.QueryResolver     { return &queryResolver{r} }
func (r *Resolver) School() generated.SchoolResolver   { return &schoolResolver{r} }

type classResolver struct{ *Resolver }
type lectureResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type schoolResolver struct{ *Resolver }
