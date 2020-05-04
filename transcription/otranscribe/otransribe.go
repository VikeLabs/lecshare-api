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

type OTR struct {
	Text      []Word `json:"text"`
	Media     string `json:"media"`
	MediaTime int    `json:"media-time"`
}

type Word struct {
	Word      string
	Timestamp float64
}

const t = `{{ range . }}<span class="wordstamp" data-timestamp="{{ printf "%.1f" .Timestamp }}">{{ .Word }}</span>
{{ end }}`

func marshalWords(words *[]Word) string {
	var buf bytes.Buffer
	t := template.Must(template.New("otr-text").Parse(t))
	err := t.Execute(&buf, *words)
	if err != nil {
		log.Panicln(err)
	}
	return buf.String()
}

func unmarshalWords(r io.Reader, data *[]Word) error {
	var word Word
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
					if attr.Key == "data-timestamp" {
						word.Timestamp, _ = strconv.ParseFloat(attr.Val, 32)
					}
				}
				tokenType = tokenizer.Next()
				if tokenType == html.TextToken {
					word.Word = tokenizer.Token().Data
					*data = append(*data, word)
				}
			}
		}

	}
	return nil
}

func (s *OTR) UnmarshalJSON(data []byte) error {
	type Alias OTR
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

func (s OTR) MarshalJSON() ([]byte, error) {
	type Alias OTR
	return json.Marshal(&struct {
		Text string `json:"text"`
		*Alias
	}{
		Text:  marshalWords(&s.Text),
		Alias: (*Alias)(&s),
	})
}
