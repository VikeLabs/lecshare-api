package graph

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"path"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/h2non/filetype"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"github.com/vikelabs/lecshare-api/graph/model"
)

// CreateLecture creates a new lecture instance within a class.
func (r *Repository) CreateLecture(ctx context.Context, input model.NewLecture, schoolCode string, courseCode string, classCode string) (*model.Lecture, error) {
	// setup
	db := r.DynamoDB
	table := db.Table(*r.TableName)
	uploader := s3manager.NewUploader(r.Session)

	// parse out the subject, code, term, section for the objectCode.
	subjectCode := strings.Split(courseCode, "#")
	termSection := strings.Split(classCode, "#")

	fmt.Println(subjectCode, termSection)

	objectKey := strings.Join([]string{schoolCode, subjectCode[0], subjectCode[1], termSection[0], termSection[1], "lectures", input.File.Filename}, "/")

	// TODO input validation

	ext := path.Ext(input.File.Filename)
	outfile := input.File.Filename[0 : len(input.File.Filename)-len(ext)]

	audioKey := strings.Join([]string{schoolCode, subjectCode[0], subjectCode[1], termSection[0], termSection[1], "lectures", outfile + ".ogg"}, "/")
	transcriptionKey := strings.Join([]string{schoolCode, subjectCode[0], subjectCode[1], termSection[0], termSection[1], "lectures", outfile + ".json"}, "/")

	// initialize blank lecture and populate
	lecture := []model.Lecture{
		{
			Name:         input.Name,
			Description:  input.Description,
			DateCreated:  time.Now(),
			DateModified: time.Now(),

			ObjectKey:     &objectKey,
			Audio:         &audioKey,
			Transcription: &transcriptionKey,
			// lecture specifics
		},
	}

	// create a new lecture entry in the table.
	err := table.Update("PK", schoolCode+"#"+courseCode).Range("SK", classCode).Append("Lectures", lecture).Run()
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("unable to create a lecture entry. please try again or contact the developers")
	}

	// upload the lecture.
	// TODO change this to be async.

	head := make([]byte, 261)
	_, err = input.File.File.Read(head)
	if err != nil {
		return nil, gqlerror.Errorf("could not identify uploaded file, please try again after verifying your file.")
	}

	kind, _ := filetype.Match(head)

	uploadReader := io.MultiReader(bytes.NewReader(head), input.File.File)

	_, err = uploader.Upload(&s3manager.UploadInput{
		// TODO remove the hardcorded value
		Bucket: r.ProcessingBucketName,
		Key:    &objectKey,
		// as we pass in an io.Reader, it will be a stream uploaded (w00t)
		Body: uploadReader,
		// TODO set additional metadata about the uploaded file.,
		ContentType: &kind.MIME.Value,
	})

	// resource.Type = kind.MIME.Value
	// resource.Size = input.File.Size

	if err != nil {
		log.Panicln("an error occurred uploading file", err)
	}
	if err != nil {
		return nil, fmt.Errorf("unable to upload file, please try again")
	}
	// TODO start transcription process. (async)
	// TODO start audio encode process. (async)

	// return the newly added lecture.
	return &lecture[0], nil
}

// CreateLectureFromResource takes a generic resource and creates a lecture.
func (r *Repository) CreateLectureFromResource(ctx context.Context, schoolCode string, courseCode string, classCode string, resourceCode string) (*model.Lecture, error) {
	resource, err := r.GetResourceByKey(ctx, schoolCode, courseCode, classCode, resourceCode)
	if err != nil {
		return nil, err
	}

	// check if the resource is valid to be turn into a lecture
	s := strings.Split(resource.Type, "/")
	if s[0] != "audio" {
		return nil, fmt.Errorf("cannot take a non-audio resource into a lecture")
	}

	// create transcription / conversion processing job.

	resource.Type = "lecture"

	return nil, nil
}

func (r *Repository) ListLectures(ctx context.Context, obj *model.Class) ([]*model.Lecture, error) {
	// converts slice types (slice of values to slice of ptrs)
	if len(obj.Lectures) > 0 {
		var lecturesRef []*model.Lecture
		for i := 0; i < len(obj.Lectures); i++ {
			Code := obj.Lectures[i].Audio
			presignedURL := r.PresignedURLGenerator.GenerateURL("/"+*Code, 120*time.Minute)
			obj.Lectures[i].Audio = &presignedURL
			lecturesRef = append(lecturesRef, &obj.Lectures[i])
		}

		return lecturesRef, nil
	}
	return nil, nil
}
