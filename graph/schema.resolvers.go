package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"time"

	"github.com/vikelabs/lecshare-api/graph/generated"
	"github.com/vikelabs/lecshare-api/graph/model"
)

func (r *classResolver) Lectures(ctx context.Context, obj *model.Class) ([]*model.Lecture, error) {
	return r.Repository.ListLectures(ctx, obj)
}

func (r *classResolver) Resources(ctx context.Context, obj *model.Class, dateBefore *time.Time, dateAfter *time.Time) ([]*model.Resource, error) {
	if dateAfter != nil || dateBefore != nil {
		return r.Repository.ListResourcesByTime(ctx, obj, dateBefore, dateAfter)
	}
	return r.Repository.ListResources(ctx, obj)
}

func (r *courseResolver) Classes(ctx context.Context, obj *model.Course, term *string) ([]*model.Class, error) {
	if term != nil {
		return r.Repository.ListClassesByTerm(ctx, obj, term)
	}
	return r.Repository.ListClasses(ctx, obj)
}

func (r *lectureResolver) Transcription(ctx context.Context, obj *model.Lecture) (*model.Transcription, error) {
	return r.Repository.GetTranscription(ctx, obj)
}

func (r *mutationResolver) CreateSchool(ctx context.Context, input model.NewSchool) (*model.School, error) {
	return r.Repository.CreateSchool(ctx, input)
}

func (r *mutationResolver) UpdateSchool(ctx context.Context, input model.UpdateSchool, schoolCode string) (*model.School, error) {
	return r.Repository.UpdateSchool(ctx, input, schoolCode)
}

func (r *mutationResolver) CreateCourse(ctx context.Context, input model.NewCourse, schoolCode string) (*model.Course, error) {
	return r.Repository.CreateCourse(ctx, input, schoolCode)
}

func (r *mutationResolver) ImportCourse(ctx context.Context, schoolCode string, courseCode string, term string) (*model.Course, error) {
	return r.Repository.ImportCourse(ctx, schoolCode, courseCode, term)
}

func (r *mutationResolver) CreateClass(ctx context.Context, input model.NewClass, schoolCode string, courseCode string) (*model.Class, error) {
	return r.Repository.CreateClass(ctx, input, schoolCode, courseCode)
}

func (r *mutationResolver) CreateResource(ctx context.Context, input model.NewResource, schoolCode string, courseCode string, classCode string) (*model.Resource, error) {
	return r.Repository.CreateResource(ctx, input, schoolCode, courseCode, classCode)
}

func (r *mutationResolver) UpdateResource(ctx context.Context, input model.UpdateResource, schoolCode string, courseCode string, classCode string, resourceKey string) (*model.Resource, error) {
	return r.Repository.UpdateResource(ctx, input, schoolCode, courseCode, classCode, resourceKey)
}

func (r *mutationResolver) CreateLecture(ctx context.Context, input model.NewLecture, schoolCode string, courseCode string, classCode string) (*model.Lecture, error) {
	return r.Repository.CreateLecture(ctx, input, schoolCode, courseCode, classCode)
}

func (r *queryResolver) Schools(ctx context.Context, code *string) ([]*model.School, error) {
	return r.Repository.ListSchools(ctx, code)
}

func (r *resourceResolver) URL(ctx context.Context, obj *model.Resource) (string, error) {
	if len(*r.Repository.CDN) > 0 {
		return r.Repository.PresignedURLGenerator.GenerateURL("/"+obj.ObjectKey, 60*time.Minute), nil
	}
	return "", nil
}

func (r *schoolResolver) Courses(ctx context.Context, obj *model.School, subject *string, code *string) ([]*model.Course, error) {
	return r.Repository.ListCourses(ctx, obj)
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

// Resource returns generated.ResourceResolver implementation.
func (r *Resolver) Resource() generated.ResourceResolver { return &resourceResolver{r} }

// School returns generated.SchoolResolver implementation.
func (r *Resolver) School() generated.SchoolResolver { return &schoolResolver{r} }

type classResolver struct{ *Resolver }
type courseResolver struct{ *Resolver }
type lectureResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type resourceResolver struct{ *Resolver }
type schoolResolver struct{ *Resolver }
