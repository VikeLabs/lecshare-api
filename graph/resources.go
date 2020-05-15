package graph

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/guregu/dynamo"
	"github.com/h2non/filetype/types"
	"github.com/rs/xid"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"github.com/vikelabs/lecshare-api/graph/model"
)

// CreateResource creates a resource record
func (r *Repository) CreateResource(ctx context.Context, input model.NewResource, schoolCode string, courseCode string, classCode string) (*model.Resource, error) {
	db := r.DynamoDB
	table := db.Table(*r.TableName)

	// input validation
	// err := r.Validate.Struct(input)
	// if err != nil {
	// 	for _, err := range err.(validator.ValidationErrors) {
	// 		graphql.AddErrorf(ctx, "field: %s, error: %s", err.StructField(), err.Tag())
	// 	}
	// 	return nil, gqlerror.Errorf("Error input errors")
	// }

	guid := xid.New()

	pk := strings.Join([]string{schoolKey, courseKey, classKey}, "#")
	sk := guid.String()

	resource := model.Resource{
		PK:           pk,
		SK:           sk,
		Name:         input.Name,
		Description:  input.Description,
		DateCreated:  time.Now(),
		DateModified: time.Now(),
	}

	if input.File != nil {
		uploader := s3manager.NewUploader(r.Session)
		var kind types.Type
		uploaderReader, err := FileTypeReader(input.File.File, &kind)

		_, err = uploader.Upload(&s3manager.UploadInput{
			// TODO remove the hardcorded value
			Bucket: r.AssetsBucketName,
			Key:    &sk,
			// as we pass in an io.Reader, it will be a stream uploaded (w00t)
			Body: uploaderReader,
			// TODO set additional metadata about the uploaded file.,
			Metadata: aws.StringMap(map[string]string{
				"Parent-Key": pk,
			}),
			ContentType: &kind.MIME.Value,
		})

		resource.ContentType = kind.MIME.Value
		resource.Type = "file"
		resource.Size = input.File.Size

		if err != nil {
			log.Panicln("an error occurred uploading file", err)
		}
	}

	// attempt to put into table if it does not exist.
	err := table.Put(resource).If("attribute_not_exists(PK)").Run()
	if err != nil {
		return nil, gqlerror.Errorf("Error: enable to create new Resource record.")

	}
	// return newly created school instance.
	return &resource, nil
}

func (r *Repository) UpdateResource(ctx context.Context, input model.UpdateResource, schoolKey string, courseKey string, classKey string, resourceKey string) (*model.Resource, error) {
	db := r.DynamoDB
	table := db.Table(*r.TableName)

	resource := model.Resource{}

	pk := strings.Join([]string{schoolKey, courseKey, classKey}, "#")

	err := table.Get("PK", pk).Range("SK", dynamo.Equal, resourceKey).One(&resource)
	if err != nil {
		return nil, gqlerror.Errorf("unable to find specified resource.")
	}

	if input.File != nil {
		uploader := s3manager.NewUploader(r.Session)
		var kind types.Type
		uploaderReader, err := FileTypeReader(input.File.File, &kind)

		_, err = uploader.Upload(&s3manager.UploadInput{
			// TODO remove the hardcorded value
			Bucket: r.AssetsBucketName,
			Key:    &resourceKey,
			// as we pass in an io.Reader, it will be a stream uploaded (w00t)
			Body: uploaderReader,
			// TODO set additional metadata about the uploaded file.,
			Metadata: aws.StringMap(map[string]string{
				"Parent-Key": pk,
			}),
			ContentType: &kind.MIME.Value,
		})

		resource.ContentType = kind.MIME.Value
		resource.Type = "file"
		resource.Filename = input.File.Filename
		resource.Size = input.File.Size
		// resource.ObjectKey = objectKey

		if err != nil {
			log.Panicln("an error occurred uploading file", err)
		}
	}

	if input.Name != nil {
		resource.Name = input.Name
	}

	if input.Description != nil {
		resource.Description = input.Description
	}

	if input.Published != nil {
		resource.Published = input.Published
	}

	resource.DateModified = time.Now()

	err = table.Put(resource).Run()
	if err != nil {
		return nil, gqlerror.Errorf("Error: enable to update Resource record.")

	}

	return &resource, nil
}

func (r *Repository) ListResources(ctx context.Context, obj *model.Class) ([]*model.Resource, error) {
	db := r.DynamoDB
	table := db.Table(*r.TableName)

	// we can only unmarshal data into a slice (not a ptr slice).
	var resources []model.Resource
	// since it is expected to return a slice of ptrs, we make a slice of ptrs.
	var resourcesRef []*model.Resource

	pk := strings.Join([]string{obj.PK, obj.Term, obj.Section}, "#")
	table.Get("PK", pk).All(&resources)

	// convert to slice of pointers.
	for i := 0; i < len(resources); i++ {
		resourcesRef = append(resourcesRef, &resources[i])
	}

	return resourcesRef, nil
}

func (r *Repository) ListResourcesByTime(ctx context.Context, obj *model.Class, dateBefore *time.Time, dateAfter *time.Time) ([]*model.Resource, error) {
	db := r.DynamoDB
	table := db.Table(*r.TableName)

	// we can only unmarshal data into a slice (not a ptr slice).
	var resources []model.Resource
	// since it is expected to return a slice of ptrs, we make a slice of ptrs.
	var resourcesRef []*model.Resource
	var before, after string

	pk := strings.Join([]string{obj.PK, obj.Term, obj.Section}, "#")
	if dateBefore != nil {
		before = xid.NewWithTime(*dateBefore).String()
	}

	if dateAfter != nil {
		after = xid.NewWithTime(*dateAfter).String()
	}

	if dateBefore != nil && dateAfter != nil {
		table.Get("PK", pk).Range("SK", dynamo.Between, before, after).All(&resources)
	} else if dateBefore != nil {
		table.Get("PK", pk).Range("SK", dynamo.LessOrEqual, before).All(&resources)
	} else if dateAfter != nil {
		table.Get("PK", pk).Range("SK", dynamo.GreaterOrEqual, after).All(&resources)
	} else {
		log.Panicln("undefined behaviour")
	}

	// convert to slice of pointers.
	for i := 0; i < len(resources); i++ {
		resourcesRef = append(resourcesRef, &resources[i])
	}

	return resourcesRef, nil
}

func (r *Repository) GetResourceByKey(ctx context.Context, schoolCode string, courseCode string, classCode string, resourceKey string) (*model.Resource, error) {
	table := r.DynamoDB.Table(*r.TableName)

	pk := strings.Join([]string{schoolCode, courseCode, classCode}, "#")
	var resource model.Resource

	err := table.Get("PK", pk).Range("SK", dynamo.Equal, resourceKey).One(&resource)
	if err != nil {
		return nil, err
	}

	return &resource, nil
}
