package dynamodb

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	dynamo "github.com/aws/aws-sdk-go/service/dynamodb"
)

// delete is an internal function responsible for remove an item of the DynamoDB table using the settings
// specified in your configuration file. If the item is removed successfully, you will receive a 200 (Ok)
// status code in response of your request. However, if something goes wrong, you will receive a 500 status
// code and an error specifying the problem
//
// you need to send the values of the keys in your request to properly remove the item
//
// To use this function, you need to specify the 'TableName' and 'Keys' in your configuration file.
func (c *DynamoDBClient) delete() (events.APIGatewayProxyResponse, error) {
	keys, err := conf.KeyCondition()
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, fmt.Errorf("failed load key conditions config: %v", err)
	}

	deleteInput := dynamo.DeleteItemInput{
		TableName: aws.String(conf.Resources.Database.TableName),
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
