package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var webhookURL string

type SlackRequestBody struct {
	Text string `json:"text"`
}

// Credits to https://golangcode.com/send-slack-messages-without-a-library/
func sendSlackNotification(webhookURL string, msg string) error {
	slackBody, _ := json.Marshal(SlackRequestBody{Text: msg})
	req, err := http.NewRequest(http.MethodPost, webhookURL, bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	if buf.String() != "ok" {
		return errors.New("Non-ok response returned from Slack")
	}
	return nil
}

func handler(ctx context.Context, event events.SNSEvent) error {
	var wg sync.WaitGroup
	wg.Add(len(event.Records))
	for _, r := range event.Records {
		go func(msg string) {
			defer wg.Done()
			sendSlackNotification(webhookURL, msg)
		}(r.SNS.Message)
	}
	wg.Wait()
	return nil
}

func init() {
	webhookURL = os.Getenv("SLACK_WEBHOOK_URL")
}

func main() {
	lambda.Start(handler)
}
