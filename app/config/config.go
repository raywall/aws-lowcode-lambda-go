package config

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	dynamo "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"gopkg.in/yaml.v2"
)

var Global = Config{}

type (
	Config struct {
		TemplateFormatVersion string    `yaml:"TemplateFormatVersion,omitempty"`
		Description           string    `yaml:"Description,omitempty"`
		Resources             Resources `yaml:"Resources"`
	}

	ComposedData struct {
		Data     string `yaml:"Data,omitempty"`
		Operator string `yaml:"Operator,omitempty"`
	}

	Database struct {
		TableName      string                  `yaml:"TableName"`
		Keys           map[string]ComposedData `yaml:"Keys"`
		Filter         string                  `yaml:"Filter,omitempty"`
		FilterValues   map[string]string       `yaml:"FilterValues,omitempty"`
		ProjectionCols []string                `yaml:"ProjectionCols,omitempty"`
	}

	KeyCondition struct {
		ExpressionAttributeNames  map[string]*string                `json:"expressionAttributeNames,omitempty"`
		ExpressionAttributeValues map[string]*dynamo.AttributeValue `json:"expressionAttributeValues,omitempty"`
		PrimaryKeys               map[string]*dynamo.AttributeValue `json:"primaryKeys,omitempty"`
		Condition                 string                            `json:"condition,omitempty"`
	}

	Resources struct {
		Request  Request  `yaml:"Request"`
		Database Database `yaml:"Database,omitempty"`
		Response Response `yaml:"Response,omitempty"`
	}

	Request struct {
		AllowedMethods []string          `yaml:"AllowedMethods"`
		HTTPMethod     string            `yaml:"HttpMethod,omitempty"`
		Parameters     map[string]string `yaml:"Parameters,omitempty"`
	}

	Response struct {
		DataStruct string `yaml:"DataStruct,omitempty"`
	}
)

func (config *Config) Load(data []byte) error {
	err := yaml.Unmarshal(data, config)
	if err != nil {
		return err
	}

	return nil
}

func (config *Config) Set() error {
	if config != nil {
		Global = *config
		return nil
	}

	return errors.New("configuration cannot be null")
}

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

func (config *Config) IsMethodAllowed() bool {
	for _, item := range config.Resources.Request.AllowedMethods {
		if item == config.Resources.Request.HTTPMethod {
			return true
		}
	}

	return false
}
