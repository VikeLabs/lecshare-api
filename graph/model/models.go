package model

import "time"

type Lecture struct {
	Title         *string `json:"title"`
	Description   *string `json:"description"`
	Datetime      *string `json:"datetime"`
	Audio         *string `json:"audio"`
	Duration      int     `json:"duration"`
	Transcription *string `json:"transcription"`
}

// School is the model used by GraphQL and DynamoDB
// this is msnaully updated.
type School struct {
	// DynamoDB
	PK string `json:"id"`
	SK string `json:"sk"`
	// Attributes
	Name         string    `json:"name" dynamo:",omitempty"`
	Code         string    `json:"code" dynamo:",omitempty"`
	Description  *string   `json:"description" dynamo:",omitempty"`
	Homepage     *string   `json:"homepage" dynamo:",omitempty"`
	Classes      []*Class  `json:"classes" dynamo:"-"`
	DateCreated  time.Time `json:"dateCreated" dynamo:",omitempty"`
	DateModified time.Time `json:"dateModified" dynamo:",omitempty"`
}
