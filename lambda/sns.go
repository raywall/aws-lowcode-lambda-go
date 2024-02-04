package lambda

import "github.com/aws/aws-lambda-go/events"

func handleSNSEvent(event events.SNSEvent) string {
	return "SNS event received"
}
