package lambda

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func serializeLocalRequest(src map[string]interface{}, dest interface{}) error {
	if src == nil || dest == nil {
		return errors.New("src and dest cannot be nil")
	}

	jsonData, err := json.Marshal(src)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonData, dest)
	if err != nil {
		return err
	}

	return nil
}

// func serializeAvro(data map[string]interface{}, schemaPath string) ([]byte, error) {
// 	avroSchema, err := os.ReadFile(schemaPath)
// 	if err != nil {
// 		return nil, err
// 	}

// 	codec, err := goavro.NewCodec(string(avroSchema))
// 	if err != nil {
// 		return nil, err
// 	}

// 	output, err := codec.TextualFromNative(nil, data)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return output, nil
// }

func marshalAttributeNames(raw map[string]interface{}) (map[string]*string, error) {
	data := make(map[string]*string)
	for key := range raw {
		if _, ok := conf.Resources.Database.Keys[key]; ok {
			data[fmt.Sprintf("#%s", key)] = aws.String(key)
		}
	}

	return data, nil
}

func marshalAttributeValues(raw map[string]interface{}) (map[string]*dynamodb.AttributeValue, error) {
	data := make(map[string]interface{})
	for key, value := range raw {
		if _, ok := conf.Resources.Database.Keys[key]; ok {
			data[fmt.Sprintf(":%s", key)] = value
		}
	}

	result, err := dynamodbattribute.MarshalMap(data)
	if err != nil {
		return nil, err
	}

	return result, nil
}
