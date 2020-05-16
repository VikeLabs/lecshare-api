package graph

import (
	"context"
	"fmt"
	"time"

	"github.com/guregu/dynamo"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"github.com/vikelabs/lecshare-api/graph/model"
)

// ListAllClasses lists all classes.
func (r *Repository) ListClasses(ctx context.Context, obj *model.Course) ([]*model.Class, error) {
	// setup
	db := r.DynamoDB
	table := db.Table(*r.TableName)

	// auth := ctx.Value("ACCESS").(string)

	var classes []model.Class
	var classesRef []*model.Class

	table.Get("PK", obj.PK+"#"+obj.SK).All(&classes)

	if len(classes) != 0 {
		for i := 0; i < len(classes); i++ {
			// c := classes[i]
			// if len(c.AccessKey) > 0 && c.AccessKey != auth {
			// 	continue
			// }
			classesRef = append(classesRef, &classes[i])
		}

		return classesRef, nil
	}
	return nil, nil
}

func (r *Repository) ListClassesByTerm(ctx context.Context, obj *model.Course, term string) ([]*model.Class, error) {
	db := r.DynamoDB
	table := db.Table(*r.TableName)

	var classes []model.Class
	var classesRef []*model.Class

	table.Get("PK", obj.PK+"#"+obj.SK).Range("SK", dynamo.BeginsWith, term).All(&classes)

	if len(classes) != 0 {
		for i := 0; i < len(classes); i++ {
			classesRef = append(classesRef, &classes[i])
		}

		return classesRef, nil
	}
	return nil, nil
}

// CreateClass creates a new instance of a class from a course in the database.
func (r *Repository) CreateClass(ctx context.Context, input model.NewClass, schoolKey string, courseKey string) (*model.Class, error) {
	// setup DynamoDB
	db := r.DynamoDB
	table := db.Table(*r.TableName)

	// TODO validate schoolKey, courseKey
	// TODO validate input

	// Grab course info.
	var course model.Course
	err := table.Get("PK", schoolKey).Range("SK", dynamo.Equal, courseKey).One(&course)
	if err != nil {
		return nil, fmt.Errorf("unable to find a valid course")
	}

	// create new class instance
	class := model.Class{
		PK:           schoolKey + "#" + courseKey,
		SK:           input.Term + "#" + input.Section,
		Name:         course.Name,
		Subject:      course.Subject,
		Code:         course.Code,
		Term:         input.Term,
		Section:      input.Section,
		DateCreated:  time.Now(),
		DateModified: time.Now(),
	}

	// insert into table.
	err = table.Put(class).If("attribute_not_exists(PK)").Run()
	if err != nil {
		return nil, gqlerror.Errorf("Error: enable to create new Class record.")

	}
	return &class, nil
}
