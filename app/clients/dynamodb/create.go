package dynamodb

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	dynamo "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func (c *DynamoDBClient) create(input *RequestInput, data *map[string]interface{}) (events.APIGatewayProxyResponse, error) {
	item, err := dynamodbattribute.MarshalMap(*data)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
	}

	putInput := &dynamo.PutItemInput{
		TableName: aws.String(input.TableName),
		Item:      item,
	}

	_, err = c.svc.PutItem(putInput)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, fmt.Errorf("falha ocorrida ao registrar os dados na tabela: %v", err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 201,
	}, nil
}
