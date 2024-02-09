package lowcode

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/raywall/aws-lowcode-lambda-go/config"
	"github.com/raywall/aws-lowcode-lambda-go/lowcodeattribute"
	"github.com/raywall/aws-lowcode-lambda-go/receiver"
)

type Lowcode interface {
	NewWithConfig(filePath string) error
	HandleRequest(ctx context.Context, evt interface{}) (interface{}, error)
}

type Function config.Settings

func (function *Function) NewWithConfig(filePath string) error {
	function.Config = &config.Config{}

	awsConfig := aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))}
	awsConfig.Endpoint = aws.String(os.Getenv("DYNAMO_ENDPOINT"))

	sess, _ := session.NewSession(&awsConfig)
	function.Client = dynamodb.New(sess)

	// read a configuration file content
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("failed reading lowcode role file: %v", err)
	}

	// load configuration
	err = function.Config.Load(data)
	if err != nil {
		log.Fatalf("failed loading settings: %v", err)
	}

	return nil
}

func (function *Function) HandleRequest(ctx context.Context, evt interface{}) (interface{}, error) {
	var event interface{} = evt

	log.Printf("received type: %T", event)

	if _, sam := os.LookupEnv("DYNAMO_ENDPOINT"); sam {
		obj := events.APIGatewayProxyRequest{}
		lowcodeattribute.SerializeLocalRequest(evt.(map[string]interface{}), &obj)
		event = obj
	}

	settings := receiver.Settings{
		Config: function.Config,
		Client: function.Client,
	}

	switch e := event.(type) {
	case events.APIGatewayProxyRequest:
		var api receiver.ApiGateway = &settings
		return api.HandleAPIGatewayEvent(e).ToGatewayResponse()
	case events.SNSEvent:
		var sns receiver.SNS = &settings
		return sns.HandleSNSEvent(e), nil
	case events.SQSEvent:
		var sqs receiver.SQS = &settings
		return sqs.HandleSQSEvent(e), nil
	case events.DynamoDBEvent:
		var dynamo receiver.DynamoDB = &settings
		return dynamo.HandleDynamoDBEvent(e), nil
	default:
		return "", fmt.Errorf("event unsupported: %T", e)
	}
}
