package graph

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/guregu/dynamo"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"github.com/vikelabs/lecshare-api/graph/model"
	"github.com/vikelabs/lecshare-api/utils"
)

// CreateCourse creates a new course entity in the database.
func (r *Repository) CreateCourse(ctx context.Context, input model.NewCourse, schoolCode string) (*model.Course, error) {
	// setup DynamoDB
	db := r.DynamoDB
	table := db.Table(*r.TableName)

	// TODO validate schoolCode

	// create new course instance
	course := model.Course{
		PK:          schoolCode,
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

func (r *Repository) ImportCourse(ctx context.Context, schoolCode string, courseCode string, term string) (*model.Course, error) {
	jsonFile, err := os.Open("/home/aomi/lecshare-api/graph/uvic_courses_kuali.json")
	if err != nil {
		log.Panicln("unable to read file")
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)

	defer jsonFile.Close()

	var courses utils.Course

	json.Unmarshal(byteValue, &courses)

	// remove all #'s
	courseCode = strings.ReplaceAll(courseCode, "#", "")
	for i := 0; i < len(courses); i++ {
		c := courses[i]
		log.Println(c)
		if c.CatalogCourseID == courseCode {
			course := model.Course{
				Name: c.Title,
			}
			return &course, nil
		}
	}
	return nil, gqlerror.Errorf("unable to import course")
}

// ListCourses lists all courses.
func (r *Repository) ListCourses(ctx context.Context, obj *model.School) ([]*model.Course, error) {
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

func (r Repository) ListCoursesBySubject(ctx context.Context, schoolCode string, subjectCode string) ([]*model.Course, error) {
	db := r.DynamoDB
	table := db.Table(*r.TableName)

	var courses []model.Course
	var coursesRef []*model.Course

	err := table.Get("PK", schoolCode).Range("SK", dynamo.BeginsWith, subjectCode).AllWithContext(ctx, &courses)
	if err != nil {
		return nil, err
	}

	// convert to slice of pointers.
	for i := 0; i < len(courses); i++ {
		coursesRef = append(coursesRef, &courses[i])
	}
	return coursesRef, nil
}
