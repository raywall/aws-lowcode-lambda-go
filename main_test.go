package main

import (
	"os"
	"testing"

	"github.com/raywall/aws-lowcode-lambda-go/config"
	"github.com/raywall/aws-lowcode-lambda-go/server/clients/dynamodb"
	"github.com/raywall/aws-lowcode-lambda-go/server/handlers"
)

func TestLoadConfig(t *testing.T) {
	conf := &config.Global

	data, err := os.ReadFile("sample.yaml")
	if err != nil {
		t.Errorf("failed reading config sample file: %v", err)
	}

	if got := conf.Load(data); got != nil {
		t.Errorf("failed loading config sample data: %v", err)
	}

	if handlers.Client, err = dynamodb.NewDynamoDBClient(conf); err != nil {
		t.Errorf("failed creating a new DynamoDBClient: %v", err)
	}
}
