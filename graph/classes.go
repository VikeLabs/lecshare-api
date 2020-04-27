package graph

import (
	"context"
	"fmt"
	"time"

	"github.com/guregu/dynamo"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"github.com/vikelabs/lecshare-api/graph/model"
)

func (r *mutationResolver) CreateClass(ctx context.Context, input model.NewClass, schoolKey string, courseKey string) (*model.Class, error) {
	// setup DynamoDB
	db := r.DB
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

func (r *courseResolver) Classes(ctx context.Context, obj *model.Course) ([]*model.Class, error) {
	// setup
	db := r.DB
	table := db.Table(*r.TableName)

	var classes []model.Class
	var classesRef []*model.Class

	table.Get("PK", obj.PK+"#"+obj.SK).All(&classes)

	// convert to slice of pointers.
	for i := 0; i < len(classes); i++ {
		classesRef = append(classesRef, &classes[i])
	}

	return classesRef, nil
}
