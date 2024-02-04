package lambda

import (
	"github.com/aws/aws-lambda-go/events"
)

func handleDynamoDBEvent(event events.DynamoDBEvent) string {
	return "DynamoDB event received"
}
