package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/vikelabs/lecshare-api/graph/generated"
	"github.com/vikelabs/lecshare-api/graph/model"
)

func (r *classResolver) Lectures(ctx context.Context, obj *model.Class) ([]*model.Lecture, error) {
	return r.Repository.ListAllLectures(ctx, obj)
}

func (r *courseResolver) Classes(ctx context.Context, obj *model.Course) ([]*model.Class, error) {
	return r.Repository.ListAllClasses(ctx, obj)
}

func (r *lectureResolver) Transcription(ctx context.Context, obj *model.Lecture) (*model.Transcription, error) {
	return r.Repository.GetTranscription(ctx, obj)
}

func (r *mutationResolver) CreateSchool(ctx context.Context, input model.NewSchool) (*model.School, error) {
	return r.Repository.CreateSchool(ctx, input)
}

func (r *mutationResolver) CreateCourse(ctx context.Context, input model.NewCourse, schoolKey string) (*model.Course, error) {
	return r.Repository.CreateCourse(ctx, input, schoolKey)
}

func (r *mutationResolver) CreateClass(ctx context.Context, input model.NewClass, schoolKey string, courseKey string) (*model.Class, error) {
	return r.Repository.CreateClass(ctx, input, schoolKey, courseKey)
}

func (r *mutationResolver) CreateLecture(ctx context.Context, input model.NewLecture, schoolKey string, courseKey string, classKey string) (*model.Lecture, error) {
	return r.Repository.CreateLecture(ctx, input, schoolKey, courseKey, classKey)
}

func (r *queryResolver) Schools(ctx context.Context, code *string) ([]*model.School, error) {
	return r.Repository.ListAllSchools(ctx, code)
}

func (r *schoolResolver) Courses(ctx context.Context, obj *model.School) ([]*model.Course, error) {
	return r.Repository.ListAllCourses(ctx, obj)
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
