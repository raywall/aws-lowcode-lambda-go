package dynamodb

import (
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	dynamo "github.com/aws/aws-sdk-go/service/dynamodb"
)

func (c *DynamoDBClient) update(data *map[string]interface{}) (events.APIGatewayProxyResponse, error) {
	keys, err := conf.KeyCondition()
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, fmt.Errorf("failed load key conditions config: %v", err)
	}

	updateInput := &dynamo.UpdateItemInput{
		TableName:                 aws.String(conf.Resources.Database.TableName),
		Key:                       keys.PrimaryKeys,
		ExpressionAttributeNames:  map[string]*string{},
		ExpressionAttributeValues: map[string]*dynamo.AttributeValue{},
	}

	updateMode := "SET"

	if mode, exist := conf.Resources.Request.Parameters["mode"]; exist {
		updateMode = mode
	}

	updateExpression := []string{}
	for key, value := range *data {
		if _, ok := conf.Resources.Database.Keys[key]; !ok {
			updateExpression = append(updateExpression, fmt.Sprintf("#%s = :%s", key, key))

			updateInput.ExpressionAttributeNames[fmt.Sprintf("#%s", key)] = aws.String(key)
			updateInput.ExpressionAttributeValues[fmt.Sprintf(":%s", key)] = &dynamo.AttributeValue{
				S: aws.String(value.(string)),
			}
		}
	}

	updateInput.UpdateExpression = aws.String(fmt.Sprintf("%s %s", updateMode, strings.Join(updateExpression, ",")))

	_, err = c.svc.UpdateItem(updateInput)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, fmt.Errorf("failed on update item: %v", err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}
