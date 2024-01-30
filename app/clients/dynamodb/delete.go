package dynamodb

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	dynamo "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func (c *DynamoDBClient) delete(input *RequestInput, data *map[string]interface{}) (events.APIGatewayProxyResponse, error) {
	keys, err := dynamodbattribute.MarshalMap(input.Keys)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
	}

	deleteInput := dynamo.DeleteItemInput{
		TableName: aws.String(input.TableName),
		Key:       keys,
	}

	_, err = c.svc.DeleteItem(&deleteInput)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}
