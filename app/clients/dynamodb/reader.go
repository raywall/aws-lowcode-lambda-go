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
func (c *DynamoDBClient) query(input *RequestInput) (events.APIGatewayProxyResponse, error) {
	var conditions string

	query := &dynamo.QueryInput{
		TableName:                 aws.String(input.TableName),
		ExpressionAttributeNames:  map[string]*string{},
		ExpressionAttributeValues: map[string]*dynamo.AttributeValue{},
	}

	for key, value := range input.Keys {
		if len(conditions) > 0 {
			conditions += " AND "
		}

		conditions += fmt.Sprintf("%s = %s", key, value.Operator, key)

		query.ExpressionAttributeNames[fmt.Sprintf("#%s", key)] = aws.String(key)
		query.ExpressionAttributeValues[fmt.Sprintf(":%s", key)] = &dynamo.AttributeValue{
			S: aws.String(value.Data),
		}
	}

	query.KeyConditionExpression = aws.String(conditions)

	// solicitar colunas especificas da consulta
	if len(input.ProjectionCols) > 0 {
		projectionCols := ""

		for i, column := range input.ProjectionCols {
			if i > 0 {
				projectionCols += ","
			}

			projectionCols += fmt.Sprintf("#%s", column)
			query.ExpressionAttributeNames[fmt.Sprintf("#%s", column)] = aws.String(column)
		}

		query.ProjectionExpression = aws.String(projectionCols)
	}

	// adiciona as regras de validacao
	if input.Filter != "" {
		query.FilterExpression = aws.String(input.Filter)

		for key, value := range input.FilterValues {
			query.ExpressionAttributeNames[fmt.Sprintf("#%s", key)] = aws.String(key)
			query.ExpressionAttributeValues[fmt.Sprintf(":%s", key)] = &dynamo.AttributeValue{
				S: aws.String(value),
			}
		}
	}

	// executa a query
	result, err := c.svc.Query(query)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
	}

	// retorna o resultado em JSON
	var jsonMap []map[string]interface{}
	if input.DataStruct != "" {
		err := json.Unmarshal([]byte(input.DataStruct), &jsonMap)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
			}, err
		}

		err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &jsonMap)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
			}, fmt.Errorf("failed to unmarshal record: %v", err)
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
			}, fmt.Errorf("erro ao serializar o resultado mapeado: %v", err)
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
		}, fmt.Errorf("erro ao serializar o resultado: %v", err)
	}

	return events.APIGatewayProxyResponse{
		Body:       string(response),
		StatusCode: 200,
	}, nil
}
