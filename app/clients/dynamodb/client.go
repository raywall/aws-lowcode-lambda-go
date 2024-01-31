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

// DynamoDBClient estrutura para o cliente do DynamoDB
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

type ComposedData struct {
	Data     string `json:"data"`
	Operator string `json:"operator"`
}

// RequestInput estrutura para a entrada de dados da consulta
type RequestInput struct {
	HTTPMethod     string                  `json:"httpMethod"`
	AllowedMethods []string                `json:"allowedMethods"`
	TableName      string                  `json:"tableName"`
	Keys           map[string]ComposedData `json:"keys"`
	Filter         string                  `json:"filter,omitempty"`
	FilterValues   map[string]string       `json:"filterValues,omitempty"`
	ProjectionCols []string                `json:"projectionCols,omitempty"`
	DataStruct     string                  `json:"dataStruct,omitempty"`
}

func (c *DynamoDBClient) HandleRequest(input *RequestInput, data *map[string]interface{}) (events.APIGatewayProxyResponse, error) {
	if !input.allowed {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
		}, fmt.Errorf("method not supported: %s", input.HTTPMethod)
	}

	switch input.HTTPMethod {
	case "GET":
		return c.query(input)
	case "POST":
		return c.create(input, data)
	case "PUT":
		return c.update(input, data)
	case "DELETE":
		return c.delete(input, data)
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
		}, fmt.Errorf("method not supported: %s", input.HTTPMethod)
	}
}

func (req *RequestInput) allowed(target string) bool {
	for _, item := range req.AllowedMethods {
		if item == target {
			return true
		}
	}

	return false
}
