package lambda

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
	"github.com/raywall/aws-lowcode-lambda-go/lambda/attributes"
	"github.com/raywall/aws-lowcode-lambda-go/lambda/resources"
)

type LowcodeFunction struct {
	Settings config.Config
	Client   *dynamodb.DynamoDB
}

func (function *LowcodeFunction) FromConfigFile(filePath string, resource string, destination string) error {
	conf := &config.Global

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
	err = conf.Load(data)
	if err != nil {
		log.Fatalf("failed loading settings: %v", err)
	}

	return nil
}

func (function *LowcodeFunction) HandleRequest(ctx context.Context, evt interface{}) (interface{}, error) {
	var event interface{} = evt

	if _, sam := os.LookupEnv("DYNAMO_ENDPOINT"); sam {
		obj := events.APIGatewayProxyRequest{}
		attributes.SerializeLocalRequest(evt.(map[string]interface{}), &obj)
		event = obj
	}

	switch e := event.(type) {
	case events.APIGatewayProxyRequest:
		return resources.HandleAPIGatewayEvent(e, function.Client).ToGatewayResponse()
	case events.SNSEvent:
		return resources.HandleSNSEvent(e, function.Client), nil
	case events.SQSEvent:
		return resources.HandleSQSEvent(e, function.Client), nil
	case events.DynamoDBEvent:
		return resources.HandleDynamoDBEvent(e, function.Client), nil
	default:
		return "", fmt.Errorf("event unsupported: %T", e)
	}
}
