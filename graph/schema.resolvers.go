package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/vikelabs/lecshare-api/graph/generated"
	"github.com/vikelabs/lecshare-api/graph/model"
)

func (r *classResolver) Lectures(ctx context.Context, obj *model.Class) ([]*model.Lecture, error) {
	return obj.Lectures, nil
}

func (r *lectureResolver) Transcription(ctx context.Context, obj *model.Lecture) (*model.Transcription, error) {
	transcriptionFile := obj.Transcription
	buff := &aws.WriteAtBuffer{}

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})

	fmt.Println(*obj.Transcription)

	downloader := s3manager.NewDownloader(sess)
	numBytes, err := downloader.Download(buff,
		&s3.GetObjectInput{
			Bucket: aws.String("assets-lecshare.oimo.ca"),
			Key:    aws.String(*transcriptionFile),
		})

	fmt.Println("Downloaded", *transcriptionFile, numBytes, "bytes")

	var transcription model.Transcription

	err = json.Unmarshal(buff.Bytes(), &transcription.Sections)
	if err != nil {
		fmt.Println(err)
	}

	return &transcription, nil
}

func (r *queryResolver) Schools(ctx context.Context, code *string) ([]*model.School, error) {
	// open schools.json
	schoolsJSONFile, err := os.Open("public/schools.json")
	if err != nil {
		panic(err)
	}

	defer schoolsJSONFile.Close()

	// read into a byte[]
	byteValue, _ := ioutil.ReadAll(schoolsJSONFile)

	var schools []*model.School
	var schoolsFiltered []*model.School

	json.Unmarshal(byteValue, &schools)

	if code != nil && len(*code) > 0 {
		for _, v := range schools {
			if v.Code == *code {
				schoolsFiltered = append(schoolsFiltered, v)
			}
		}
		return schoolsFiltered, nil
	}
	return schools, nil
}

func (r *schoolResolver) Classes(ctx context.Context, obj *model.School) ([]*model.Class, error) {
	classesJSONFile, err := os.Open("public/classes.json")
	if err != nil {
		panic(err)
	}

	defer classesJSONFile.Close()

	byteValue, _ := ioutil.ReadAll(classesJSONFile)

	var classes map[string][]*model.Class
	json.Unmarshal(byteValue, &classes)

	fmt.Println(obj.Code)

	c := classes[strings.ToLower(obj.Code)]
	if c != nil {
		return c, nil
	}

	return nil, nil
}

// Class returns generated.ClassResolver implementation.
func (r *Resolver) Class() generated.ClassResolver { return &classResolver{r} }

// Lecture returns generated.LectureResolver implementation.
func (r *Resolver) Lecture() generated.LectureResolver { return &lectureResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// School returns generated.SchoolResolver implementation.
func (r *Resolver) School() generated.SchoolResolver { return &schoolResolver{r} }

type classResolver struct{ *Resolver }
type lectureResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type schoolResolver struct{ *Resolver }
