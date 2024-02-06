package attributes

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/raywall/aws-lowcode-lambda-go/config"
)

// Conf refers to the lambda function's configuration, containing all the necessary information
// about the request, database, and response parameters that the client uses to orchestrate requests.
var conf = &config.Global

func SerializeLocalRequest(src map[string]interface{}, dest interface{}) error {
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

func MarshalAttributeNames(raw map[string]interface{}, separator string) (map[string]*string, error) {
	data := make(map[string]*string)
	for key := range raw {
		if _, ok := conf.Resources.Database.Keys[key]; ok {
			data[fmt.Sprintf("%s%s", separator, key)] = aws.String(key)
		}
	}

	return data, nil
}

func MarshalAttributeValues(raw map[string]interface{}, separator string) (map[string]*dynamodb.AttributeValue, error) {
	data := make(map[string]interface{})
	for key, value := range raw {
		if _, ok := conf.Resources.Database.Keys[key]; ok {
			data[fmt.Sprintf("%s%s", separator, key)] = value
		}
	}

	result, err := dynamodbattribute.MarshalMap(data)
	if err != nil {
		return nil, err
	}

	return result, nil
}
