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
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/guregu/dynamo"
	"github.com/vektah/gqlparser/v2/gqlerror"
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

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println("Downloaded", *transcriptionFile, numBytes, "bytes")

	var transcription model.Transcription

	err = json.Unmarshal(buff.Bytes(), &transcription.Sections)
	if err != nil {
		fmt.Println(err)
	}

	return &transcription, nil
}

func (r *mutationResolver) CreateSchool(ctx context.Context, input model.NewSchool) (*model.School, error) {
	// note: the same table is used used accross the entire base application.
	tableName := os.Getenv("tableName")
	fmt.Println(tableName)
	// TODO get session from ctx
	db := dynamo.New(session.New(), &aws.Config{Region: aws.String("us-west-2")})
	table := db.Table(tableName)

	school := model.School{
		PK:           "ORG",
		SK:           input.Code,
		Name:         input.Name,
		Code:         input.Code,
		Description:  input.Description,
		Homepage:     input.Homepage,
		DateCreated:  time.Now(),
		DateModified: time.Now(),
	}

	err := table.Put(school).If("attribute_not_exists(PK)").Run()
	if err != nil {
		return nil, gqlerror.Errorf("Error: enable to create new School record.")

	}
	return &school, nil
}

func (r *queryResolver) Schools(ctx context.Context, code *string) ([]*model.School, error) {
	// note: the same table is used used accross the entire base application.
	tableName := os.Getenv("tableName")
	// TODO get session from ctx
	db := dynamo.New(session.New(), &aws.Config{Region: aws.String("us-west-2")})
	table := db.Table(tableName)

	var schools []model.School
	var schoolsRef []*model.School
	table.Get("PK", "ORG").All(&schools)

	// convert to slice of pointers.
	for i := 0; i < len(schools); i++ {
		schoolsRef = append(schoolsRef, &schools[i])
	}

	return schoolsRef, nil
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

	// Swap all audio for the Vikelabs one
	for i := range classes {
		for _, c := range classes[i] {
			for _, L := range c.Lectures {
				*L.Audio, err = GetResource("vikelabs/vikelabs_test1.ogg", 15)
				if err != nil {
					println(err)
					return nil, err
				}
			}
		}
	}

	c := classes[strings.ToLower(obj.Code)]
	if c != nil {
		return c, nil
	}

	return nil, nil
}

func (r *schoolResolver) DateCreated(ctx context.Context, obj *model.School) (string, error) {
	return obj.DateCreated.String(), nil
}

func (r *schoolResolver) DateModified(ctx context.Context, obj *model.School) (string, error) {
	return obj.DateCreated.String(), nil
}

// Class returns generated.ClassResolver implementation.
func (r *Resolver) Class() generated.ClassResolver { return &classResolver{r} }

// Lecture returns generated.LectureResolver implementation.
func (r *Resolver) Lecture() generated.LectureResolver { return &lectureResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// School returns generated.SchoolResolver implementation.
func (r *Resolver) School() generated.SchoolResolver { return &schoolResolver{r} }

type classResolver struct{ *Resolver }
type lectureResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type schoolResolver struct{ *Resolver }
