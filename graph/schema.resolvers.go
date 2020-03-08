// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
package graph

import (
	"context"

	"github.com/vikelabs/lecshare-api/graph/generated"
	"github.com/vikelabs/lecshare-api/graph/model"
)

func (r *lectureResolver) Transcription(ctx context.Context, obj *model.Lecture) (*model.Transcription, error) {
	return getTranscription("public/vikelabs_test1.json")
}

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

func (r *Resolver) Lecture() generated.LectureResolver { return &lectureResolver{r} }
func (r *Resolver) Query() generated.QueryResolver     { return &queryResolver{r} }

type lectureResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
