package dynamodb

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/raywall/aws-lowcode-lambda-go/config"
)

// Client refers to the DynamoDBClient created for integrating the Lambda function with DynamoDB.
// It enables the redirection of GET, POST, PUT, and DELETE requests.
var Client *DynamoDBClient

// HandleLambdaEvent is the function responsible for receiving API Gateway requests in your Lambda function.
// It redirects these requests to the DynamoDBClient, which is responsible for interacting with the database
// and returning a response with the result.
//
// Before using this function, you'll need to build a Config object with request, database, and response
// configuration, and instantiate the DynamoDBClient.
//
// If everything is okay, you will receive an APIGatewayProxyResponse. However, if something goes wrong,
// you'll receive an empty response and an error describing the problem.
func HandleAPIGatewayEvent(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var jsonMap map[string]interface{}
	conf := &config.Global

	err := json.Unmarshal([]byte(event.Body), &jsonMap)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, fmt.Errorf("failed to unmarshal request body: %v", err)
	}

	conf.Resources.Request.HTTPMethod = event.HTTPMethod
	conf.Resources.Request.Parameters = event.QueryStringParameters

	for key, reg := range conf.Resources.Database.Keys {
		reg.Data = jsonMap[key].(string)
		conf.Resources.Database.Keys[key] = reg
	}

	return Client.handleRequest(&jsonMap)
}

func HandleDynamoDBStreamEvent(ctx context.Context, event events.DynamoDBEvent) (events.DynamoDBEventResponse, error) {

	if len(event.Records) > 0 {

		for _, item := range event.Records {
			switch item.EventName {
			case "MODIFY":

			case "REMOVE":

			default:
			}
		}
	}

	return events.DynamoDBEventResponse{}, nil
}

func HandleSQSEvent(ctx context.Context, event events.DynamoDBEvent) (events.DynamoDBEventResponse, error) {
	response := events.DynamoDBEventResponse{}

	if len(event.Records) > 0 {
		for _, evt := range event.Records {
			switch evt.EventName {
			case "MODIFY":
				evt.
			case "REMOVE":

			default:
			}
		}
	}

	return response, nil
}
