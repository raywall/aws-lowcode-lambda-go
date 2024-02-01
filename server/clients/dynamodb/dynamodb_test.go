package dynamodb

import (
	"os"
	"testing"

	"github.com/raywall/aws-lowcode-lambda-go/config"
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

	if Client, err = NewDynamoDBClient(conf); err != nil {
		t.Errorf("failed creating a new DynamoDBClient: %v", err)
	}
}
