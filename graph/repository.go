package graph

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-playground/validator/v10"
	"github.com/guregu/dynamo"
)

// Repository is _
type Repository struct {
	// Data Validation
	Validate *validator.Validate
	// AWS DynamoDB
	DynamoDB  *dynamo.DB
	TableName *string
	// AWS S3
	Session              *session.Session
	ProcessingBucketName *string
	AssetsBucketName     *string
}
