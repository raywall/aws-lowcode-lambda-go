package dynamodb

import (
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	dynamo "github.com/aws/aws-sdk-go/service/dynamodb"
)

// update is an internal function responsible for update, remove or add the attributes of an item
// into the DynamoDB table specified in your configuration file. If the item is updated successfully, you
// will receive a 200 (Ok) status code in response of your request. However, if something goes wrong, you
// will receive a 500 status code and an error specifying the problem.
//
// the content with the item attributes nedds to be in the body of the request to proceed with the update
//
// To use this function, you need to specify the 'TableName' and 'Keys' in your configuration file.
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
