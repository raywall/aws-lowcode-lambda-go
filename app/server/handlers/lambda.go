package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/raywall/aws-lowcode-lambda-go/config"
	"github.com/raywall/aws-lowcode-lambda-go/server/clients/dynamodb"

	"github.com/aws/aws-lambda-go/events"
)

var Client *dynamodb.DynamoDBClient

func HandleLambdaEvent(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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

	return Client.HandleRequest(&jsonMap)
}
