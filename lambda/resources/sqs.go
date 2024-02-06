package resources

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func HandleSQSEvent(event events.SQSEvent, client *dynamodb.DynamoDB) string {
	return "SQS event received"
}
