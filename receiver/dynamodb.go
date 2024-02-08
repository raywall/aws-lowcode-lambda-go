package receiver

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func HandleDynamoDBEvent(event events.DynamoDBEvent, client *dynamodb.DynamoDB) string {
	return "DynamoDB event received"
}
