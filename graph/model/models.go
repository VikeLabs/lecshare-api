package model

type Lecture struct {
	Title         *string `json:"title"`
	Description   *string `json:"description"`
	Datetime      *string `json:"datetime"`
	Audio         *string `json:"audio"`
	Duration      int     `json:"duration"`
	Transcription *string `json:"transcription"`
}
