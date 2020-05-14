package graph

import (
	"bytes"
	"io"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-playground/validator/v10"
	"github.com/guregu/dynamo"
	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
	"github.com/vikelabs/lecshare-api/utils/bunnycdn"
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

	CDN                   *string
	PresignedURLGenerator *bunnycdn.Generator
}

// FileTypeReader will identity the MIME type from an io.Reader
func FileTypeReader(r io.Reader, t *types.Type) (io.Reader, error) {
	head := make([]byte, 261)
	_, err := r.Read(head)
	if err != nil {
		return nil, err
	}

	*t, err = filetype.Match(head)
	if err != nil {
		return nil, err
	}
	return io.MultiReader(bytes.NewReader(head), r), nil
}
