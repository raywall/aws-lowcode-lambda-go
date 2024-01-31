package dynamodb

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	dynamo "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// query realiza uma consulta no DynamoDB com base no input
func (c *DynamoDBClient) query() (events.APIGatewayProxyResponse, error) {
	var queryInput = &dynamodb.QueryInput{
		TableName:                 aws.String(input.TableName),
		ExpressionAttributeNames:  map[string]*string{},
		ExpressionAttributeValues: map[string]*dynamo.AttributeValue{},
	}

	keyCondition, err := conf.KeyCondition()
	if err != nil {
		return event.APIGatewayProxyResponse{
			StatusCode: 500,
		}, fmt.Errorf("failed to load key conditions: %v", err)
	}

	queryInput.KeyConditionExpression = aws.String(keyCondition.Condition)
	queryInput.ExpressionAttributeNames = keyCondition.ExpressionAttributeNames
	queryInput.ExpressionAttributeValues = keyCondition.ExpressionAttributeValues

	// get list of specific cols
	if len(conf.Database.ProjectionCols) > 0 {
		projectionCols := ""

		for i, col := range conf.Database.ProjectionCols {
			if i > 0 {
				projectionCols += ","
			}

			projectionCols += fmt.Sprintf("#%s ", col)
			queryInput.ExpressionAttributeNames[fmt.Sprintf("#%s", col)] = aws.String(col)
		}

		queryInput.ProjectionCols = aws.String(projectionCols)
	}

	// add validation rules
	if conf.Database.Filter != "" {
		queryInput.FilterExpression = aws.String(conf.Database.Filter)

		for key, value := range conf.Database.FilterValues {
			queryInput.ExpressionAttributeNames[fmt.Sprintf("#%s", key)] = aws.String(key)
			queryInput.ExpressionAttributeValues[fmt.Sprintf(":%s", key)] = &dynamo.AttributeValue{
				S: aws.String(value),
			}
		}
	}

	// run database query
	result, err := c.svc.Query(queryInput)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, fmt.Errorf("failed to run database query: %v", err)
	}

	// return a JSON based response
	var jsonMap []map[string]interface{}
	if conf.Response.DataStruct != "" {
		err := json.Unmarshal([]byte(conf.Response.DataStruct), &jsonMap)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
			}, fmt.Errorf("failed unmarshal data struct config: %v", err)
		}

		err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &jsonMap)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
			}, fmt.Errorf("failed unmarshal record: %v", err)
		}

		var data interface{}
		switch len(jsonMap) {
		case 0:
			data = nil
		case 1:
			data = jsonMap[0]
		default:
			data = jsonMap
		}

		jsonResponse, err := json.Marshal(data)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
			}, fmt.Errorf("failed marshal mapped response: %v", err)
		}

		return events.APIGatewayProxyResponse{
			Body:       string(jsonResponse),
			StatusCode: 200,
		}, nil
	}

	response, err := json.Marshal(result.Items)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, fmt.Errorf("failed marshal response: %v", err)
	}

	return events.APIGatewayProxyResponse{
		Body:       string(response),
		StatusCode: 200,
	}, nil
}
