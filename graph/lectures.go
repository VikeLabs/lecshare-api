package graph

import (
	"context"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/vikelabs/lecshare-api/graph/model"
)

// CreateLecture creates a new lecture instance within a class.
func (r *Repository) CreateLecture(ctx context.Context, input model.NewLecture, schoolKey string, courseKey string, classKey string) (*model.Lecture, error) {
	// setup
	db := r.DynamoDB
	table := db.Table(*r.TableName)
	uploader := s3manager.NewUploader(r.Session)

	// parse out the subject, code, term, section for the objectkey.
	subjectCode := strings.Split(courseKey, "#")
	termSection := strings.Split(classKey, "#")

	fmt.Println(subjectCode, termSection)

	objectKey := strings.Join([]string{schoolKey, subjectCode[0], subjectCode[1], termSection[0], termSection[1], "lectures", input.File.Filename}, "/")

	// TODO input validation

	ext := path.Ext(input.File.Filename)
	outfile := input.File.Filename[0 : len(input.File.Filename)-len(ext)]

	audioKey := strings.Join([]string{schoolKey, subjectCode[0], subjectCode[1], termSection[0], termSection[1], "lectures", outfile + ".ogg"}, "/")
	transcriptionKey := strings.Join([]string{schoolKey, subjectCode[0], subjectCode[1], termSection[0], termSection[1], "lectures", outfile + ".json"}, "/")

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
	err := table.Update("PK", schoolKey+"#"+courseKey).Range("SK", classKey).Append("Lectures", lecture).Run()
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("unable to create a lecture entry. please try again or contact the developers")
	}

	// upload the lecture.
	// TODO change this to be async.
	_, err = uploader.Upload(&s3manager.UploadInput{
		// TODO remove the hardcorded value
		Bucket: r.ProcessingBucketName,
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

func (r *Repository) ListAllLectures(ctx context.Context, obj *model.Class) ([]*model.Lecture, error) {
	// converts slice types (slice of values to slice of ptrs)
	// TODO sign with BunnyCDN pre-signed for security / CDN.
	fmt.Println(obj)
	fmt.Println(len(obj.Lectures))
	if len(obj.Lectures) > 0 {

		svc := s3.New(r.Session)
		var lecturesRef []*model.Lecture
		for i := 0; i < len(obj.Lectures); i++ {
			key := obj.Lectures[i].Audio
			req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
				Bucket: r.AssetsBucketName,
				Key:    key,
			})
			presignedURL, _ := req.Presign(60 * time.Minute)
			obj.Lectures[i].Audio = &presignedURL
			lecturesRef = append(lecturesRef, &obj.Lectures[i])
		}

		return lecturesRef, nil
	}
	return nil, nil
}
