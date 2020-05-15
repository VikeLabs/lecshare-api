package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var webhookURL string 

// Credits to https://golangcode.com/send-slack-messages-without-a-library/
func sendSlackNotification(webhookURL string, msg string) error {
    slackBody, _ := json.Marshal(SlackRequestBody{Text: msg})
    req, err := http.NewRequest(http.MethodPost, webhookUrl, bytes.NewBuffer(slackBody))
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
	for _, r := range event.Records {
		sendSlackNotification(webhookURL, r.SNS.Message)
	}
	return nil
}

func init(){
	webhookURL = os.Getenv("SLACK_WEBHOOK_URL")
}

func main() {
	lambda.Start(handler)
}