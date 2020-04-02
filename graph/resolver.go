// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.
package graph

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-playground/validator/v10"
	"github.com/guregu/dynamo"
)

type Resolver struct {
	Session    *session.Session
	DB         *dynamo.DB
	BucketName *string
	TableName  *string
	Validate   *validator.Validate
}
