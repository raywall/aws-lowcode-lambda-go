package dynamodb

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	dynamo "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func (c *DynamoDBClient) update(input *RequestInput, data *map[string]interface{}) (events.APIGatewayProxyResponse, error) {
	item, err := dynamodbattribute.MarshalMap(data)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
	}

	updateInput := &dynamo.UpdateItemInput{
		TableName: aws.String(input.TableName),
		Key:       item,
	}

	_, err = c.svc.UpdateItem(updateInput)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, fmt.Errorf("falha ocorrida ao atualizar os dados dor egistro na tabela: %v", err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}
