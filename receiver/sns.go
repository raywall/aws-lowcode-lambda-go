package receiver

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func HandleSNSEvent(event events.SNSEvent, client *dynamodb.DynamoDB) string {
	return "SNS event received"
}
