package receiver

import (
	"github.com/aws/aws-lambda-go/events"
)

type DynamoDB interface {
	HandleDynamoDBEvent(event events.DynamoDBEvent) string
}

func (s *Settings) HandleDynamoDBEvent(event events.DynamoDBEvent) string {
	return "DynamoDB event received"
}
