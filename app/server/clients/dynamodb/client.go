package dynamodb

import (
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	dynamo "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var conf = Config.Global

type DynamoDBClient struct {
	svc dynamodbiface.DynamoDBAPI
}

// NewDynamoDBClient cria uma nova instancia do cliente do DynamoDB
func NewDynamoDBClient() (*DynamoDBClient, error) {
	config := aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))}

	if endpoint, present := os.LookupEnv("ENDPOINT"); present {
		config.Endpoint = aws.String(endpoint)
	}

	sess, err := session.NewSession(&config)
	if err != nil {
		return &DynamoDBClient{}, fmt.Errorf("erro ao iniciar sessão aws: %v", err)
	}
	return &DynamoDBClient{
		svc: dynamo.New(sess),
	}, nil
}

func (c *DynamoDBClient) HandleRequest(data *map[string]interface{}) (events.APIGatewayProxyResponse, error) {
	if !conf.allowed {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
		}, fmt.Errorf("method not supported: %s", conf.HTTPMethod)
	}

	switch conf.HTTPMethod {
	case "GET":
		return c.query()
	case "POST":
		return c.create(data)
	case "PUT":
		return c.update(data)
	case "DELETE":
		return c.delete(data)
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
		}, fmt.Errorf("method not supported: %s", HTTPMethod)
	}
}

func (config *Config) allowed(target string) bool {
	for _, item := range config.AllowedMethods {
		if item == target {
			return true
		}
	}

	return false
}
