// This Go package provides ready-to-use handlers designed for quick and easy usage without
// the need for extensive code writing, often referred to as 'low-code.'
// These handlers are constructed based on the origin of the request being made to the Lambda function.
// The package's documentation covers the public API, including details about available handlers, their
// signatures, accepted input types, and produced output types. Users can benefit from clear usage examples,
// code comments that explain the logic behind each handler, installation and configuration instructions,
// error handling guidance, and any external dependencies required. Additionally, the documentation showcases
// real-world use cases, integration examples with other services or frameworks, and recommendations for
// testing to ensure the handlers function as intended. Users can find licensing information, attribution
// requirements, and details on how to stay updated with package releases and contribute to its development.
// The package aims to simplify the development process by providing developers with pre-built handlers that
// can be easily integrated into their AWS Lambda functions.
package dynamodb

import (
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	dynamo "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/raywall/aws-lowcode-lambda-go/config"
)

// Conf refers to the lambda function's configuration, containing all the necessary information
// about the request, database, and response parameters that the client uses to orchestrate requests.
var conf = &config.Global

// DynamoDBClient is a struct representing the DynamoDB client. It includes the 'svc' attribute,
// which is an interface to enable mocking the dynamodb.DynamoDB service client's API operation,
// paginators, and waiters. This make unit testing your code that calls out to the SDK's service
// client's calls easier.
type DynamoDBClient struct {
	svc dynamodbiface.DynamoDBAPI
}

// NewDynamoDBClient is responsible for creating a new DynamoDBClient, which is a crucial function.
//
// If you have a 'DYNAMO_ENDPOINT' environment variable pointing to a local DynamoDB container defined,
// the session will automatically be configured to use this endpoint.
//
// This function creates a session and loads the configurations if provided, returning a pointer
// reference to the DynamoDBClient.
//
// If something goes wrong, you will receive an empty client and an error.
func NewDynamoDBClient(configuration ...*config.Config) (*DynamoDBClient, error) {
	config := aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))}

	if endpoint, present := os.LookupEnv("DYNAMO_ENDPOINT"); present {
		config.Endpoint = aws.String(endpoint)
	}

	sess, err := session.NewSession(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to start an aws session: %v", err)
	}

	client := &DynamoDBClient{
		svc: dynamo.New(sess),
	}

	if len(configuration) > 0 {
		err = configuration[0].Set()
		if err != nil {
			return nil, fmt.Errorf("failed to load configuration: %v", err)
		}
	}

	return client, nil
}

// handleRequest is an internal BFF (Backend For Frontend) request responsible for checking if the HTTP method
// in the request is allowed and calling actions based on the HTTP method.
//
// If you receive a request that isn't supported, you will receive a 400 error status code.
func (c *DynamoDBClient) handleRequest(data *map[string]interface{}) (events.APIGatewayProxyResponse, error) {

	if !conf.IsMethodAllowed() {
		return events.APIGatewayProxyResponse{
			StatusCode: 401,
			Body:       fmt.Sprintf("%s method is not allowed", conf.Resources.Request.HTTPMethod),
		}, nil
	}

	switch conf.Resources.Request.HTTPMethod {
	case "GET":
		return c.query()
	case "POST":
		return c.create(data)
	case "PUT":
		return c.update(data)
	case "DELETE":
		return c.delete()
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf("%s method is not supported", conf.Resources.Request.HTTPMethod),
		}, nil
	}
}
