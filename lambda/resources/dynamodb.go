package resources

import (
	"github.com/aws/aws-lambda-go/events"
)

func HandleDynamoDBEvent(event events.DynamoDBEvent) string {
	return "DynamoDB event received"
}
