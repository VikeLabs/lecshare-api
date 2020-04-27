package graph

import (
	"context"
	"fmt"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/go-playground/validator"
	"github.com/guregu/dynamo"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"github.com/vikelabs/lecshare-api/graph/model"
)

func (r *mutationResolver) CreateSchool(ctx context.Context, input model.NewSchool) (*model.School, error) {
	// TODO get session from ctx
	db := r.DB
	table := db.Table(*r.TableName)

	// input validation
	err := r.Validate.Struct(input)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			graphql.AddErrorf(ctx, "field: %s, error: %s", err.StructField(), err.Tag())
		}
		return nil, gqlerror.Errorf("Error input errors")
	}

	// create new school instance
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

	// attempt to put into table if it does not exist.
	err = table.Put(school).If("attribute_not_exists(PK)").Run()
	if err != nil {
		return nil, gqlerror.Errorf("Error: enable to create new School record.")

	}
	// return newly created school instance.
	return &school, nil
}

func (r *queryResolver) Schools(ctx context.Context, code *string) ([]*model.School, error) {
	// note: the same table is used used accross the entire base application.
	db := r.DB
	table := db.Table(*r.TableName)

	// we can only unmarshal data into a slice (not a ptr slice).
	var schools []model.School
	// since it is expected to return a slice of ptrs, we make a slice of ptrs.
	var schoolsRef []*model.School

	// if given a code, filter to only return that school.
	if code != nil {
		fmt.Printf("<%s>\n", *code)
		schools = make([]model.School, 1)
		schoolsRef = make([]*model.School, 1)
		err := table.Get("PK", "ORG").Range("SK", dynamo.Equal, *code).One(&schools[0])
		fmt.Printf("%v\n", schools[0])
		if err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("unable to find scholl with code %s", *code)
		}
		schoolsRef[0] = &schools[0]
		return schoolsRef, nil
	}

	table.Get("PK", "ORG").All(&schools)

	// convert to slice of pointers.
	for i := 0; i < len(schools); i++ {
		schoolsRef = append(schoolsRef, &schools[i])
	}

	return schoolsRef, nil
}
