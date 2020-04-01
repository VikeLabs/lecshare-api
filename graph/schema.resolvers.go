package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/guregu/dynamo"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"github.com/vikelabs/lecshare-api/graph/generated"
	"github.com/vikelabs/lecshare-api/graph/model"
)

func (r *classResolver) Lectures(ctx context.Context, obj *model.Class) ([]*model.Lecture, error) {
	// converts slice types (slice of values to slice of ptrs)
	// TODO sign with BunnyCDN pre-signed for security / CDN.

	var lecturesRef []*model.Lecture
	for i := 0; i < len(obj.Lectures); i++ {
		lecturesRef = append(lecturesRef, &obj.Lectures[i])
	}

	return lecturesRef, nil
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

func (r *lectureResolver) Transcription(ctx context.Context, obj *model.Lecture) (*model.Transcription, error) {
	// when it retrieves the file, it stores it in memory (this buffer below) rather than on disk.
	buff := &aws.WriteAtBuffer{}

	downloader := s3manager.NewDownloader(r.Session)
	_, err := downloader.Download(buff,
		&s3.GetObjectInput{
			Bucket: r.BucketName,
			Key:    obj.Transcription,
		})

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// fmt.Println("Downloaded", *transcriptionFile, numBytes, "bytes")

	var transcription model.Transcription

	err = json.Unmarshal(buff.Bytes(), &transcription.Sections)
	if err != nil {
		fmt.Println(err)
	}

	return &transcription, nil
}

func (r *mutationResolver) CreateSchool(ctx context.Context, input model.NewSchool) (*model.School, error) {
	// TODO get session from ctx
	db := r.DB
	table := db.Table(*r.TableName)

	// TODO input validation

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
	err := table.Put(school).If("attribute_not_exists(PK)").Run()
	if err != nil {
		return nil, gqlerror.Errorf("Error: enable to create new School record.")

	}
	// return newly created school instance.
	return &school, nil
}

func (r *mutationResolver) CreateCourse(ctx context.Context, input model.NewCourse, schoolKey string) (*model.Course, error) {
	// setup DynamoDB
	db := r.DB
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

func (r *mutationResolver) CreateLecture(ctx context.Context, input model.NewLecture, schoolKey string, courseKey string, classKey string) (*model.Lecture, error) {
	// setup
	db := r.DB
	table := db.Table(*r.TableName)
	uploader := s3manager.NewUploader(r.Session)

	// parse out the subject, code, term, section for the objectkey.
	subjectCode := strings.Split(courseKey, "#")
	termSection := strings.Split(classKey, "#")
	objectKey := strings.Join([]string{schoolKey, subjectCode[0], subjectCode[1], termSection[0], termSection[1], "lectures", input.File.Filename}, "/")

	// TODO input validation

	// initialize blank lecture and populate
	lecture := []model.Lecture{
		model.Lecture{
			Name:         input.Name,
			Description:  input.Description,
			DateCreated:  time.Now(),
			DateModified: time.Now(),
			ObjectKey:    &objectKey,
			Audio:        aws.String(""),
			// lecture specifics
		},
	}

	// create a new lecture entry in the table.
	err := table.Update("PK", schoolKey+"#"+courseKey).Range("SK", classKey).Append("Lectures", lecture).Run()
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("unable to create a lecture entry. please try again or contact the developers")
	}

	// upload the lecture.
	// TODO change this to be async.
	_, err = uploader.Upload(&s3manager.UploadInput{
		// TODO remove the hardcorded value
		Bucket: r.BucketName,
		Key:    aws.String(objectKey),
		// as we pass in an io.Reader, it will be a stream uploaded (w00t)
		Body: input.File.File,
		// TODO set additional metadata about the uploaded file.
	})
	if err != nil {
		return nil, fmt.Errorf("unable to upload file, please try again")
	}

	// TODO start transcription process. (async)
	// TODO start audio encode process. (async)

	// return the newly added lecture.
	return &lecture[0], nil
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

func (r *schoolResolver) Courses(ctx context.Context, obj *model.School) ([]*model.Course, error) {
	db := r.DB
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
