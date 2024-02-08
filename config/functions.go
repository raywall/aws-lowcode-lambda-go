package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	schema "github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v2"
)

// carega as configuracoes do arquivo de configuração
func (c *Config) FromFile(path string) error {
	// read a configuration file content
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// load configuration
	err = yaml.Unmarshal(data, c)
	if err != nil {
		return err
	}

	return nil
}

// Valida formato do json de acordo com o schema
func (res *ResourceItem) Validate(data map[string]interface{}) bool {
	content, err := os.ReadFile(res.ObjectPathSchema)
	if err != nil {
		fmt.Printf("failed open config file: %v", err)
		return false
	}

	schemaLoader := schema.NewStringLoader(string(content))
	documentLoader := schema.NewGoLoader(data)

	result, err := schema.Validate(schemaLoader, documentLoader)
	if err != nil {
		fmt.Println("Erro durante a validação:", err)
		return false
	}

	return result.Valid()
}

// Pega o schema do arquivo json
func (res *ResourceItem) UnmarshalSchema() (map[string]interface{}, error) {
	model, err := os.ReadFile(res.ObjectPathSchema)
	if err != nil {
		return nil, err
	}

	var jsonMap map[string]interface{}
	json.Unmarshal(model, &jsonMap)

	return jsonMap, nil
}

// Estrutura dados para criação de registro
func (res *ResourceItem) MarshalMap(data map[string]interface{}) (map[string]*dynamodb.AttributeValue, error) {
	if !res.Validate(data) {
		return nil, errors.New("unsupported data structure")
	}

	response, err := dynamodbattribute.MarshalMap(data)
	if err != nil {
		return nil, err
	}

	return response, nil
}
