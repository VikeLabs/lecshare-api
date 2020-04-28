package model

type TranscriptionJSON struct {
	Results *TranscriptionResultsJSON `json:"results"`
}

type TranscriptionResultsJSON struct {
	Transcripts []*TranscriptionTranscriptJSON `json:"transcripts"`
	Items       []*TranscriptionItemJSON       `json:"items"`
}

type TranscriptionTranscriptJSON struct {
	Transcript *string `json:"transcript"`
}

type TranscriptionItemJSON struct {
	Type         string                   `json:"type"`
	StartTime    *string                  `json:"start_time"`
	EndTime      *string                  `json:"end_time"`
	Alternatives []*TranscriptionWordJSON `json:"alternatives"`
}
type TranscriptionWordJSON struct {
	Content string `json:"content"`
}
