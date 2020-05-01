package graph

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/vektah/gqlparser/v2/gqlerror"
	"github.com/vikelabs/lecshare-api/graph/model"
	"github.com/vikelabs/lecshare-api/utils"
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

func (r *Repository) ImportCourse(ctx context.Context, schoolKey string, courseKey string, term string) (*model.Course, error) {
	jsonFile, err := os.Open("/home/aomi/lecshare-api/graph/uvic_courses_kuali.json")
	if err != nil {
		log.Panicln("unable to read file")
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)

	defer jsonFile.Close()

	var courses utils.Course

	json.Unmarshal(byteValue, &courses)

	// remove all #'s
	courseKey = strings.ReplaceAll(courseKey, "#", "")
	for i := 0; i < len(courses); i++ {
		c := courses[i]
		log.Println(c)
		if c.CatalogCourseID == courseKey {
			course := model.Course{
				Name: c.Title,
			}
			return &course, nil
		}
	}
	return nil, gqlerror.Errorf("unable to import course")
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
