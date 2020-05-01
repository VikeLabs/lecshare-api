package model

import "time"

// School is the model used by GraphQL and DynamoDB
// this is manually updated.
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

type NewSchool struct {
	Name        string  `json:"name" validate:"required"`
	Code        string  `json:"code" validate:"gte=4,lte=8"`
	Description *string `json:"description"`
	Homepage    *string `json:"homepage" validate:"url"`
}

type Course struct {
	// DynamoDB
	PK string `json:"id"`
	SK string `json:"sk"`
	// Attributes
	Name        string  `json:"name"`
	Subject     string  `json:"subject"`
	Code        string  `json:"code"`
	Description *string `json:"description"`
	Homepage    *string `json:"homepage"`
}

type Class struct {
	// DynamoDB
	PK string `json:"id"`
	SK string `json:"sk"`
	// Attributes
	Name         string        `json:"name"`
	Subject      string        `json:"subject"`
	Code         string        `json:"code"`
	Term         string        `json:"term"`
	Section      string        `json:"section"`
	Instructors  []*Instructor `json:"instructors" dynamo:",omitempty"`
	Lectures     []Lecture     `json:"lectures"`
	DateCreated  time.Time     `json:"dateCreated"`
	DateModified time.Time     `json:"dateModified"`
}

type Lecture struct {
	// DynamoDB
	PK string `json:"id"`
	SK string `json:"sk"`
	// Attributes
	Name          *string   `json:"title"`
	Description   *string   `json:"description"`
	Audio         *string   `json:"audio"`
	Duration      int       `json:"duration"`
	Transcription *string   `json:"transcription"`
	ObjectKey     *string   `json:"objectKey"`
	DateCreated   time.Time `json:"dateCreated"`
	DateModified  time.Time `json:"dateModified"`
}

type Resource struct {
	// DynamoDB
	PK string `json:"id"`
	SK string `json:"sk"`
	// Attributes
	Name         *string   `json:"name"`
	Description  *string   `json:"description"`
	ObjectKey    string    `json:"objectKey"`
	Type         string    `json:"type"`
	Size         int64     `json:"size"`
	Published    *bool     `json:"published"`
	DateCreated  time.Time `json:"dateCreated"`
	DateModified time.Time `json:"dateModified"`
}
