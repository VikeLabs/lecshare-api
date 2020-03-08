package graph

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/vikelabs/lecshare-api/graph/model"
)

func getTranscription(filename string) (*model.Transcription, error) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	bytes, _ := ioutil.ReadAll(file)

	var transcription model.Transcription

	err = json.Unmarshal(bytes, &transcription.Sections)
	if err != nil {
		fmt.Println(err)
	}

	return &transcription, nil
}
