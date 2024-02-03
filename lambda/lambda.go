package lambda

import (
	"errors"
	"log"
	"os"

	"github.com/jmespath/go-jmespath"
	"github.com/raywall/aws-lowcode-lambda-go/config"
	"github.com/raywall/aws-lowcode-lambda-go/server/clients/dynamodb"
)

const _ = jmespath.ASTEmpty

type LowcodeFunction struct {
	Handler     interface{}
	Settings    config.Config
	Resource    string
	Destination string
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

	// check resource, destination and create a handler to integrate the services
	if destination == "DYNAMODB" {
		dynamodb.Client, err = dynamodb.NewDynamoDBClient(conf)
		if err != nil {
			log.Fatalf("failed starting a dynamodb client: %v", err)
		}

		switch resource {
		case "APIGATEWAY":
			function.Handler = dynamodb.HandleAPIGatewayEvent

		case "DYNAMOSTREAM":
			function.Handler = dynamodb.HandleDynamoDBStreamEvent

		case "SQS":
			function.Handler = dynamodb.HandleSQSEvent

		default:
			return errors.ErrUnsupported
		}
	} else {
		return errors.ErrUnsupported
	}

	return nil
}
