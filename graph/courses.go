package graph

import (
	"context"

	"github.com/vektah/gqlparser/v2/gqlerror"
	"github.com/vikelabs/lecshare-api/graph/model"
)

// CreateCourse creates a new course entity in the database.
func (r *Repository) CreateCourse(ctx context.Context, input model.NewCourse, schoolKey string) (*model.Course, error) {
	// setup DynamoDB
	db := r.DynamoDB
	table := db.Table(*r.TableName)

	// TODO validate schoolKey

	// create new course instance
	course := model.Course{
		PK:          schoolKey,
		SK:          input.Subject + "#" + input.Code,
		Name:        input.Name,
		Subject:     input.Subject,
		Code:        input.Code,
		Description: &input.Description,
		Homepage:    input.Homepage,
	}
	err := table.Put(course).If("attribute_not_exists(PK)").Run()
	if err != nil {
		return nil, gqlerror.Errorf("Error: enable to create new Course record.")

	}
	return &course, nil
}

// ListAllCourses lists all courses.
func (r *Repository) ListAllCourses(ctx context.Context, obj *model.School) ([]*model.Course, error) {
	db := r.DynamoDB
	table := db.Table(*r.TableName)

	var courses []model.Course
	var coursesRef []*model.Course

	table.Get("PK", obj.Code).All(&courses)

	// convert to slice of pointers.
	for i := 0; i < len(courses); i++ {
		coursesRef = append(coursesRef, &courses[i])
	}
	return coursesRef, nil
}