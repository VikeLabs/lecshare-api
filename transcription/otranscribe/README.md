# Lecshare
## oTranscribe
The following package provides encode/decode capability of `.otr` files exported from [oTranscribe](https://otranscribe.com/).

A `otr` file is simply an JSON with an key/value pair where the value is a list of `<span>` tags containing time and word metadata. So we can treat `.otr` for the most part. 

## Usage
The following does an decode/encode of a `.otr` file.
```go
file, _ := os.Open("example.otr")

otrBytes, _ := ioutil.ReadAll(file)

var o OTR

// decode to otr
err = json.Unmarshal(otrBytes, &o)
if err != nil {
    log.Fatal(err)
}

// manipulate / read as desired

// encode to otr
outBytes, _ := json.Marshal(o)

_ = ioutil.WriteFile("encoded_example.otr", outBytes, 0644)
```