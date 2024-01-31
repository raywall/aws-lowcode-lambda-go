package dynamodb

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	dynamo "github.com/aws/aws-sdk-go/service/dynamodb"
)

func (c *DynamoDBClient) delete() (events.APIGatewayProxyResponse, error) {
	keys, err := conf.KeyCondition()
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, fmt.Errorf("failed load key conditions config: %v", err)
	}

	deleteInput := dynamo.DeleteItemInput{
		TableName: aws.String(conf.Database.TableName),
		Key:       keys.PrimaryKeys,
	}

	_, err = c.svc.DeleteItem(&deleteInput)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, fmt.Errorf("failed to remove table item: %v", err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}
