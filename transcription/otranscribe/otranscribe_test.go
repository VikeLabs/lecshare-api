package otransribe

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestOTRGeneration(t *testing.T) {
	// generate()
}

func TestOTRUnmarshal(t *testing.T) {
	file, err := os.Open("./testfiles/vikelabs_test1.otr")
	if err != nil {
		log.Fatal(err)
	}

	otrBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	var o OTR

	err = json.Unmarshal(otrBytes, &o)
	if err != nil {
		log.Fatal(err)
	}

	outBytes, err := json.Marshal(o)
	if err != nil {
		log.Fatal(err)
	}
	_ = ioutil.WriteFile("test.otr", outBytes, 0644)
}
