// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
package graph

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/vikelabs/lecshare-api/graph/generated"
	"github.com/vikelabs/lecshare-api/graph/model"
)

func (r *queryResolver) Classes(ctx context.Context) ([]*model.Class, error) {
	c := []*model.Class{
		&model.Class{
			Title: "Foundations of Programming II",
			Code:  "CSC 115",
			Instructor: &model.User{
				FirstName: "Bill",
				LastName:  "Bird",
				Suffix:    "Dr",
				Role:      "Instructor",
			},
			Lectures: []*model.Lecture{
				&model.Lecture{
					Name:     "Introduction",
					Datetime: "Feb, 12, 2020",
					Duration: 3600,
				},
			},
		},
	}
	return c, nil
}

func (r *queryResolver) Transcriptions(ctx context.Context) (*model.Transcription, error) {
	file, err := os.Open("public/vikelabs_test1.json")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	bytes, _ := ioutil.ReadAll(file)

	var transcription model.Transcription

	err = json.Unmarshal(bytes, &transcription.Sections)
	if err != nil {
		fmt.Println(err)
	}

	return &transcription, nil
}

func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
