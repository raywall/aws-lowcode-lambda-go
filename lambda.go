package lambda

import (
	"errors"
	"log"
	"os"

	"github.com/jmespath/go-jmespath"
	"github.com/raywall/aws-lowcode-lambda-go/server/clients/dynamodb"
)

const _ = jmespath.ASTEmpty

type LowcodeFunction struct {
	handler     interface{}
	settings    Config
	resource    string
	destination string
}

func (function *LowcodeFunction) FromConfigFile(filePath string, resource string, destination string) error {
	conf := &Global

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

	// check resource, destination and create a handler to integrate the services
	if destination == "DYNAMODB" {
		dynamodb.Client, err = dynamodb.NewDynamoDBClient(conf)
		if err != nil {
			log.Fatalf("failed starting a dynamodb client: %v", err)
		}

		switch resource {
		case "APIGATEWAY":
			function.handler = dynamodb.HandleAPIGatewayEvent

		case "DYNAMOSTREAM":
			function.handler = dynamodb.HandleDynamoDBStreamEvent

		case "SQS":
			function.handler = dynamodb.HandleSQSEvent

		default:
			return errors.ErrUnsupported
		}
	} else {
		return errors.ErrUnsupported
	}

	return nil
}
