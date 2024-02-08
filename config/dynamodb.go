package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Recovers the primary key from the table
func (res *ResourceItem) GetPrimaryKeyAttributeValue(data interface{}) (map[string]*dynamodb.AttributeValue, error) {
	if res.ResourceType != "DynamoDB" {
		return nil, errors.New("resource is not a database")
	}

	// if !res.Validate(data) {
	// 	return nil, errors.New("unsupported data structure")
	// }

	keys := make(map[string]interface{})
	for key := range res.Properties.Keys {
		if value, ok := data.(map[string]interface{})[key]; ok {
			keys[key] = value
		}
	}

	response, err := dynamodbattribute.MarshalMap(keys)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Returns the names of all attributes
func (res *ResourceItem) GetAllAttributeNames(data interface{}) (map[string]*string, error) {
	if res.ResourceType != "DynamoDB" {
		return nil, errors.New("the resource is not a dynamodb table")
	}

	names := make(map[string]*string)
	for key := range data.(map[string]interface{}) {
		names[fmt.Sprintf("#%s", key)] = aws.String(key)
	}

	return names, nil
}

// Returns the values of all attributes
func (res *ResourceItem) GetAllAttributeValues(data interface{}) (map[string]*dynamodb.AttributeValue, error) {
	if res.ResourceType != "DynamoDB" {
		return nil, errors.New("the resource is not a dynamodb table")
	}

	temp := make(map[string]interface{})
	for key, value := range data.(map[string]interface{}) {
		temp[fmt.Sprintf(":%s", key)] = aws.String(fmt.Sprintf("%v", value))
	}

	values, err := dynamodbattribute.MarshalMap(temp)
	if err != nil {
		return nil, err
	}

	return values, nil
}

// Returns the names of all attributes, except those who make up the primary key
func (res *ResourceItem) GetAttributeNames(data interface{}) (map[string]*string, error) {
	if res.ResourceType != "DynamoDB" {
		return nil, errors.New("the resource is not a dynamodb table")
	}

	names := make(map[string]*string)
	for key := range data.(map[string]interface{}) {
		if _, ok := res.Properties.Keys[key]; !ok {
			names[fmt.Sprintf("#%s", key)] = aws.String(key)
		}
	}

	return names, nil
}

// Returns the values of all attributes except those that make up the primary key
func (res *ResourceItem) GetAttributeValues(data interface{}) (map[string]*dynamodb.AttributeValue, error) {
	if res.ResourceType != "DynamoDB" {
		return nil, errors.New("the resource is not a dynamodb table")
	}

	temp := make(map[string]interface{})
	for key, value := range data.(map[string]interface{}) {
		if _, ok := res.Properties.Keys[key]; !ok {
			temp[fmt.Sprintf(":%s", key)] = aws.String(fmt.Sprintf("%v", value))
		}
	}

	values, err := dynamodbattribute.MarshalMap(temp)
	if err != nil {
		return nil, err
	}

	return values, nil
}

func (res *ResourceItem) GetKeyAttributeNames(data interface{}) (map[string]*string, error) {
	if res.ResourceType != "DynamoDB" {
		return nil, errors.New("the resource is not a dynamodb table")
	}

	names := make(map[string]*string)
	for key := range data.(map[string]interface{}) {
		if _, ok := res.Properties.Keys[key]; ok {
			names[fmt.Sprintf("#%s", key)] = aws.String(key)
		}
	}

	return names, nil
}

func (res *ResourceItem) GetKeyAttributeValues(data interface{}) (map[string]*dynamodb.AttributeValue, error) {
	if res.ResourceType != "DynamoDB" {
		return nil, errors.New("the resource is not a dynamodb table")
	}

	temp := make(map[string]interface{})
	for key, value := range data.(map[string]interface{}) {
		if _, ok := res.Properties.Keys[key]; ok {
			temp[fmt.Sprintf(":%s", key)] = aws.String(fmt.Sprintf("%v", value))
		}
	}

	values, err := dynamodbattribute.MarshalMap(temp)
	if err != nil {
		return nil, err
	}

	return values, nil
}

func (res *ResourceItem) GetUpdateExpression(data interface{}) (string, error) {
	commands := []string{}

	for key := range data.(map[string]interface{}) {
		if _, ok := res.Properties.Keys[key]; !ok {
			commands = append(commands, fmt.Sprintf("#%s = :%s", key, key))
		}
	}

	return fmt.Sprintf("SET %s", strings.Join(commands, ",")), nil
}

func (res *ResourceItem) GetKeyConditions(data interface{}) (string, error) {
	conditions := []string{}

	for key := range data.(map[string]interface{}) {
		if _, ok := res.Properties.Keys[key]; ok {
			conditions = append(conditions, fmt.Sprintf("#%s = :%s", key, key))
		}
	}

	return strings.Join(conditions, " AND "), nil
}
