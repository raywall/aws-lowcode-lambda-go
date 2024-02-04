package lambda

import "github.com/aws/aws-lambda-go/events"

func handleSQSEvent(event events.SQSEvent) string {
	return "SQS event received"
}
