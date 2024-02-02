package lambda

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	dynamo "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"gopkg.in/yaml.v2"
)

// Global is the object that contains your configuration
var Global = Config{}

type (
	// Config is a struct representing the Configuration file, containing all the necessary information about
	// the request, database, and response parameters that the client uses to orchestrate requests
	Config struct {
		TemplateFormatVersion string    `yaml:"TemplateFormatVersion,omitempty"`
		Description           string    `yaml:"Description,omitempty"`
		Resources             Resources `yaml:"Resources"`
	}

	// ComposedData is a struct representing an AttributeValue. You will use this struct to inform
	// parameters of an key.
	//
	// Example:
	// - let's consider that you have a PK and a SK, and your PK needs to be exactly, but your SK is a number
	//   and needs to be better than 5
	//
	// In this case you will indicate:
	// > Keys:
	//   	email:
	//     		Operator: "="
	//   	id:
	//     		Operator: "="
	//
	// You don't need to inform the data, because the client will get then of your request.
	ComposedData struct {
		Data     string `yaml:"Data,omitempty"`
		Operator string `yaml:"Operator,omitempty"`
	}

	// Database is a struct representing the necessary DynamoDB information required to perform actions on your table.
	// To use a DynamoDBClient, it's necessary to specify a 'TableName' and the 'Keys' to be used.
	// If you want to apply filtering to a query, you can specify the 'Filter' and 'FilterValues' attributes.
	// When using the 'Filter', use '#name_of_the_column' to specify the column name and ':name_of_the_column'
	// to represent the value to be filtered. In 'FilterValues', provide the actual value to be filtered.
	//
	// In the following example, we have a query on the 'users' table, specifying keys 'email' and 'product' that should
	// be equal to the values indicated in the request. Additionally, the query will only retrieve items where the 'age'
	// is greater than 18.
	//
	// Example:
	// > Database:
	// 		TableName: users
	// 		Keys:
	//   		email:
	//     			Operator: "="
	//   		id:
	//     			Operator: "="
	// 		Filter: "#age > :age"
	// 		FilterValues:
	//   		age: "18"
	// 		ProjectionCols:
	// 		- email
	// 		- username
	// 		- age
	Database struct {
		TableName      string                  `yaml:"TableName"`
		Keys           map[string]ComposedData `yaml:"Keys"`
		Filter         string                  `yaml:"Filter,omitempty"`
		FilterValues   map[string]string       `yaml:"FilterValues,omitempty"`
		ProjectionCols []string                `yaml:"ProjectionCols,omitempty"`
	}

	// KeyCondition is a struct used in the client to organize the parameters of the query.
	KeyCondition struct {
		ExpressionAttributeNames  map[string]*string                `json:"expressionAttributeNames,omitempty"`
		ExpressionAttributeValues map[string]*dynamo.AttributeValue `json:"expressionAttributeValues,omitempty"`
		PrimaryKeys               map[string]*dynamo.AttributeValue `json:"primaryKeys,omitempty"`
		Condition                 string                            `json:"condition,omitempty"`
	}

	// Resources is a struct used as an attribute of the Config to group information about request,
	// database, and response configurations.
	Resources struct {
		Request  Request  `yaml:"Request"`
		Database Database `yaml:"Database,omitempty"`
		Response Response `yaml:"Response,omitempty"`
	}

	// Request is a struct representing the configuration of an HTTP request. This struct defines the HTTP methods
	// allowed to be processed by the client, the method received in the request, and all the query parameters provided.
	//
	// Example:
	// > Request:
	// 		AllowedMethods:
	// 		- GET
	// 		- POST
	// 		- PUT
	Request struct {
		AllowedMethods []string          `yaml:"AllowedMethods"`
		HTTPMethod     string            `yaml:"HttpMethod,omitempty"`
		Parameters     map[string]string `yaml:"Parameters,omitempty"`
	}

	// Response is a struct representing the configuration of the response. This struct specifies the data structure
	// of the response in a JSON-based format.
	//
	// Example:
	// > Response:
	// 		DataStruct: '[{"username": "", "age": ""}]'
	Response struct {
		DataStruct string `yaml:"DataStruct,omitempty"`
	}
)

// Load is the function responsible for unmarshaling the configuration YAML file into an object that can be
// used by the DynamoDBClient.
func (config *Config) Load(data []byte) error {
	err := yaml.Unmarshal(data, config)
	if err != nil {
		return err
	}

	return nil
}

// Set is the function that can be used to load the content of a Config struct into the Global variable.
func (config *Config) Set() error {
	if config != nil {
		Global = *config
		return nil
	}

	return errors.New("configuration cannot be null")
}

// KeyCondition is a function used to prepare data to be used as expressionNames, expressionValues, primary key
// information, and conditions to orchestrate actions in the DynamoDB table.
func (config *Config) KeyCondition() (*KeyCondition, error) {
	var (
		condition string
		names     = map[string]*string{}
		values    = map[string]*dynamo.AttributeValue{}
		keys      = map[string]string{}
	)

	for key, value := range config.Resources.Database.Keys {
		if condition != "" {
			condition += " AND "
		}

		condition += fmt.Sprintf("#%s %s :%s", key, value.Operator, key)

		names[fmt.Sprintf("#%s", key)] = aws.String(key)
		values[fmt.Sprintf(":%s", key)] = &dynamo.AttributeValue{
			S: aws.String(value.Data),
		}

		if _, ok := config.Resources.Database.Keys[key]; ok {
			keys[key] = value.Data
		}
	}

	primaryKeys, err := dynamodbattribute.MarshalMap(keys)
	if err != nil {
		return nil, fmt.Errorf("failed unmarshal primary keys: %s", err)
	}

	return &KeyCondition{
		ExpressionAttributeNames:  names,
		ExpressionAttributeValues: values,
		Condition:                 condition,
		PrimaryKeys:               primaryKeys,
	}, nil
}

// IsMethodAllowed checks if the received method is allowed to be processed by the client.
func (config *Config) IsMethodAllowed() bool {
	for _, item := range config.Resources.Request.AllowedMethods {
		if item == config.Resources.Request.HTTPMethod {
			return true
		}
	}

	return false
}
