package dynamodb

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	dynamo "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// create is an internal function responsible for inserting a new item into the DynamoDB table specified in your
// configuration file. If the new item is created successfully, you will receive a 201 (created) status code
// in response of your request. However, if something goes wrong, you will receive a 500 status code and
// an error specifying the problem.
//
// the content with the item attributes needs to be in the body of the request
//
// To use this function, you need to specify the 'TableName' and 'Keys' in your configuration file.
func (c *DynamoDBClient) create(data *map[string]interface{}) (events.APIGatewayProxyResponse, error) {
	item, err := dynamodbattribute.MarshalMap(*data)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, fmt.Errorf("failed marshal map body information")
	}

	putInput := &dynamo.PutItemInput{
		TableName: aws.String(conf.Resources.Database.TableName),
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
