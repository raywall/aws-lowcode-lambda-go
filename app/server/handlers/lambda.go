package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/raywall/aws-lowcode-lambda-go/server/clients/dynamodb"

	"github.com/aws/aws-lambda-go/events"
)

var (
	Client    *dynamodb.DynamoDBClient
	InputData *dynamodb.RequestInput
)

func HandleLambdaEvent(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var jsonMap map[string]interface{}

	err := json.Unmarshal([]byte(event.Body), &jsonMap)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, fmt.Errorf("failed to unmarshal request body: %v", err)
	}

	InputData.HTTPMethod = event.HTTPMethod

	for key, reg := range InputData.Keys {
		reg.Data = jsonMap[key].(string)
		InputData.Keys[key] = reg
	}

	return Client.HandleRequest(InputData, &jsonMap)
}
