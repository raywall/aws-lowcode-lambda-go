package lambda

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/raywall/aws-lowcode-lambda-go/config"
)

type LowcodeFunction struct {
	Settings config.Config
}

func (function *LowcodeFunction) FromConfigFile(filePath string, resource string, destination string) error {
	conf := &config.Global

	// read a configuration file content
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("failed reading lowcode role file: %v", err)
	}

	// load configuration
	err = conf.Load(data)
	if err != nil {
		log.Fatalf("failed loading settings: %v", err)
	}

	// function.HandlerRequest = handleRequest
	return nil
}

func (function *LowcodeFunction) HandleRequest(ctx context.Context, evt interface{}) (any, error) {
	var event interface{} = evt

	if _, sam := os.LookupEnv("DYNAMO_ENDPOINT"); sam {
		obj := events.APIGatewayProxyRequest{}
		serializeLocalRequest(evt.(map[string]interface{}), &obj)
		event = obj
	}

	switch e := event.(type) {
	case events.APIGatewayProxyRequest:
		return handleAPIGatewayEvent(e)
	case events.SNSEvent:
		return handleSNSEvent(e), nil
	case events.SQSEvent:
		return handleSQSEvent(e), nil
	case events.DynamoDBEvent:
		return handleDynamoDBEvent(e), nil
	default:
		return "", fmt.Errorf("event unsupported: %T", e)
	}
}
