package resources

import "github.com/aws/aws-lambda-go/events"

func HandleSNSEvent(event events.SNSEvent) string {
	return "SNS event received"
}
