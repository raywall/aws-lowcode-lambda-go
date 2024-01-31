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

var conf = &config.Global

type DynamoDBClient struct {
	svc dynamodbiface.DynamoDBAPI
}

// NewDynamoDBClient cria uma nova instancia do cliente do DynamoDB
func NewDynamoDBClient() (*DynamoDBClient, error) {
	config := aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))}

	if endpoint, present := os.LookupEnv("DYNAMO_ENDPOINT"); present {
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
