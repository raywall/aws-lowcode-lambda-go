package resources

import "github.com/aws/aws-lambda-go/events"

func HandleSQSEvent(event events.SQSEvent) string {
	return "SQS event received"
}
