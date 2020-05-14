package otransribe

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

type Transcription struct {
	Text      []Data  `json:"text"`
	Media     string  `json:"media"`
	MediaTime float32 `json:"media-time"`
}

type Data struct {
	Data      string
	Type      string
	Timestamp float64
}

const t = `{{ range . }}<span class="{{ .Type }}" data-timestamp="{{ printf "%.1f" .Timestamp }}">{{ .Data }}</span>{{ end }}`

func marshalWords(words *[]Data) string {
	var buf bytes.Buffer
	t := template.Must(template.New("otr-text").Parse(t))
	err := t.Execute(&buf, *words)
	if err != nil {
		log.Panicln(err)
	}
	return buf.String()
}

func unmarshalWords(r io.Reader, data *[]Data) error {
	var word Data
	tokenizer := html.NewTokenizer(r)
	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			err := tokenizer.Err()
			if err == io.EOF {
				//end of the file, break out of the loop
				break
			}
			log.Fatalf("error tokenizing HTML: %v", tokenizer.Err())
		}

		if tokenType == html.StartTagToken {
			token := tokenizer.Token()
			if "span" == token.Data {
				for _, attr := range token.Attr {
					if attr.Key == "class" {
						word.Type = attr.Val
					} else if attr.Key == "data-timestamp" {
						word.Timestamp, _ = strconv.ParseFloat(attr.Val, 32)
					}
				}
				tokenType = tokenizer.Next()
				if tokenType == html.TextToken {
					word.Data = tokenizer.Token().Data
					*data = append(*data, word)
				}
			}
		}

	}
	return nil
}

// UnmarshalJSON implements OTR unmarshalling
func (s *Transcription) UnmarshalJSON(data []byte) error {
	type Alias Transcription
	aux := &struct {
		Text string `json:"text"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	_ = unmarshalWords(strings.NewReader(aux.Text), &s.Text)
	return nil
}

// MarshalJSON implements marshalling for OTR
func (s Transcription) MarshalJSON() ([]byte, error) {
	type Alias Transcription
	return json.Marshal(&struct {
		Text string `json:"text"`
		*Alias
	}{
		Text:  marshalWords(&s.Text),
		Alias: (*Alias)(&s),
	})
}
