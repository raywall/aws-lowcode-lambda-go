package config

import "github.com/aws/aws-sdk-go/service/dynamodb"

type (
	ActionRequested string

	Settings struct {
		Config *Config
		Client *dynamodb.DynamoDB
	}

	Config struct {
		TemplateFormatVersion string    `yaml:"TemplateFormatVersion"`
		Description           string    `yaml:"Description"`
		Resources             Resources `yaml:"Resources"`
	}

	Resources struct {
		Receiver  ResourceItem `yaml:"Receiver"`
		Connector ResourceItem `yaml:"Connector"`
	}

	ResourceItem struct {
		ObjectPathSchema string `yaml:"ObjectPathSchema"`
		ResourceType     string `yaml:"ResourceType"`

		Properties Properties `yaml:"Properties"`
	}

	Properties struct {
		// ApiGateway Receiver
		AllowedMethods []string          `yaml:"AllowedMethods"`
		AllowedPath    map[string]string `yaml:"AllowedPath"`

		// DynamoDB Connector
		TableName     string                 `yaml:"TableName"`
		Keys          map[string]string      `yaml:"Keys"`
		Filter        []string               `yaml:"Filters"`
		FilterValues  map[string]interface{} `yaml:"FilterValues"`
		OutputColumns []string               `yaml:"OutputColumns"`
	}

	DynamoAttributes struct {
		Key                map[string]*dynamodb.AttributeValue
		KeyAttributeValues map[string]*dynamodb.AttributeValue
		KeyAttributeNames  map[string]*string
		AttributeNames     map[string]*string
		AttributeValues    map[string]*dynamodb.AttributeValue
		KeyCondition       string
		UpdateExpression   string
	}
)

const (
	Create ActionRequested = "POST"
	Read   ActionRequested = "GET"
	Update ActionRequested = "PUT"
	Delete ActionRequested = "DELETE"
)
